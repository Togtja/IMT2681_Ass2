package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"RESTGvkGitLab/caching"
	"RESTGvkGitLab/firedb"
	"RESTGvkGitLab/globals"

	"google.golang.org/api/iterator"
)

func activateWebhook(event globals.EventMsg, params []string) error {
	fmt.Println("Calling webhook with", event)
	var invoke Invocation
	invoke.Event = string(event)
	invoke.Params = params
	invoke.Time = time.Now().String()
	invByte, err := json.Marshal(invoke)
	if err != nil {
		return err
	}
	iter := firedb.Client.Collection(globals.WebhookF).Where(globals.EventF, "==", string(event)).Documents(firedb.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		//We failed to iterate through
		if err != nil {
			return err
		}

		m := doc.Data()
		var url string = fmt.Sprint(m[globals.URLF])
		_, err = http.Post(url, "application/json", bytes.NewBuffer(invByte))
		if err != nil {
			//If we fail we just move on
			continue
			/*We ignore this error because
			There are other webhooks that need to be called
			And we don't know who's fault it is*/
		}
	}
	return nil
}

//Finds current uptime
func uptime() time.Duration {
	return time.Since(globals.StartTime)
}

//Start the timer to figure out current uptime
func init() {
	globals.StartTime = time.Now()
}

func apiGetCall(w http.ResponseWriter, getReq string, auth string, v interface{}) error {
	client := &http.Client{}
	request, err := http.NewRequest("GET", getReq, nil)
	if err != nil {
		errmsg := "The HTTP request failed with error: " + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
	}
	//Set authentication if there is one
	if auth != globals.PUBLIC {
		request.Header.Set("Private-Token", auth)
	}
	resp, err := client.Do(request)
	if err != nil {
		errmsg := "The HTTP request failed with error: " + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return err
	}
	defer resp.Body.Close()
	//Some APIcall when calling for commits return a 404,
	//However, I don't want to throw that error due to 99% of them working
	//It's pointles, but the API call return empty handed
	if resp.StatusCode == 404 {
		return nil
	}
	//Invalid authentication
	if resp.StatusCode == 401 {
		http.Error(w, "Invalid Authentication", http.StatusUnauthorized)
		return errors.New("Invalid Authentication in API call")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errmsg := "The Read of the response failed with error: " + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return err
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		errmsg := "The Unmarshal failed with error: " + err.Error()
		errmsg = errmsg + "\n Possibly failed Authentication"
		http.Error(w, errmsg, http.StatusInternalServerError)
		return err
	}
	return nil
}
func findLimit(w http.ResponseWriter, r *http.Request) int64 {
	limittxt := r.FormValue(globals.LIMITP)
	limit, err := strconv.ParseInt(limittxt, 10, 32)
	if err != nil && limittxt != "" {
		errmsg := "Invalid limit, error: " + err.Error()
		http.Error(w, errmsg, http.StatusBadRequest)
	}
	if limit <= 0 {
		//Default limit
		limit = globals.DLIMIT
	}
	return limit
}
func findOffset(w http.ResponseWriter, r *http.Request) int64 {
	offsettxt := r.FormValue(globals.OFFSETP)
	offset, err := strconv.ParseInt(offsettxt, 10, 32)
	if err != nil && offsettxt != "" {
		errmsg := "Invalid offset, error: " + err.Error()
		http.Error(w, errmsg, http.StatusBadRequest)
	}
	return offset
}
func subAPICallsForCommits(projects []Project, auth string, w http.ResponseWriter) []Repo {
	query := globals.GITAPI + "projects/"
	var repos []Repo
	//We have the project now we need to find the amount of commits for each
	//project
	var wg sync.WaitGroup
	var m = &sync.Mutex{}
	for i := range projects {
		wg.Add(1)
		//Do calls in multithreading
		go func(i int) {

			for j := 0; j < globals.MAXPAGE; j++ {

				subquery := query + strconv.Itoa(projects[i].ID) + globals.GITREPO + globals.PAGEQ + strconv.Itoa(j+1)
				var commits []Commit
				err := apiGetCall(w, subquery, auth, &commits)
				if err != nil {
					//The API call has failed
					wg.Done()
					return
				}
				if len(commits) == 0 {
					break
				}
				projects[i].Commits = append(projects[i].Commits, commits...)

			}
			//Make sure we don't append at the same time

			m.Lock()
			repos = append(repos, Repo{projects[i].NamePath, len(projects[i].Commits)})
			m.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	sort.SliceStable(repos, func(i, j int) bool {
		return repos[i].Commits > repos[j].Commits
	})
	return repos
}
func subAPICallsForLang(projects []Project, auth string, w http.ResponseWriter) []string {
	query := globals.GITAPI + globals.PROJQ
	var lang []string

	//String map to find duplicates
	dupFreq := make(map[string]int)

	//We have the project now we need to find programming languages
	//project
	var wg sync.WaitGroup
	var m = &sync.Mutex{}
	for i := range projects {
		wg.Add(1)
		//Do calls in multithreading
		go func(i int) {
			subquery := query + strconv.Itoa(projects[i].ID) + globals.LANGQ
			var v interface{}
			err := apiGetCall(w, subquery, auth, &v)
			if err != nil {
				//The API call has failed
				wg.Done()
				return
			}
			data := v.(map[string]interface{})
			var language string = ""
			for k, v := range data {
				switch v.(type) {
				case float64:
					language = k
					//Make sure we don't append at the same time
					m.Lock()
					dupFreq[language]++
					//If it is the first time we seen it
					if dupFreq[language] == 1 {

						lang = append(lang, language)
					}
					m.Unlock()

				default:
					continue
				}
			}
			if language == "" {
				wg.Done()
				return
			}

			wg.Done()
		}(i)
	}
	wg.Wait()

	sort.SliceStable(lang, func(i, j int) bool {
		return dupFreq[lang[i]] > dupFreq[lang[j]]
	})

	return lang
}
func isGetRequest(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		http.Error(w, "only get method allowed", http.StatusNotImplemented)
		return false
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Expecting format .../", http.StatusBadRequest)
		return false
	}
	return true
}

//Works for commits and languages
func genericGetHandler(w http.ResponseWriter, r *http.Request, fileName string, fileDir string,
	v interface{}, auth string) ([]Project, int64, int64, bool) {

	//Find the headers
	limit := findLimit(w, r)
	offset := findOffset(w, r)
	projects := GetProjects(w, r, auth)
	if projects == nil {
		return nil, limit, offset, false
	}

	status, file := caching.ShouldFileCache(fileName, fileDir)
	defer file.Close()
	if status == globals.Error || status == globals.DirFail {
		http.Error(w, "Failed to create a file", http.StatusInternalServerError)
		return nil, limit, offset, false
	}
	if status == globals.Exist {
		//The file exist
		err := caching.ReadFile(file, &v)
		if err != nil {
			errmsg := "The Failed Reading from file with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return nil, limit, offset, false
		}
		//We have no file
	} else {
		//If we have no project file, we have no lang or commit files
		repo, okR := v.(*Repos)
		if okR {
			repo.Repos = subAPICallsForCommits(projects, auth, w)
			v = repo
		}
		lang, okL := v.(*Lang)
		if okL {
			fmt.Println("Finding languages")
			lang.Language = subAPICallsForLang(projects, auth, w)
			v = lang
		}
		if !okR && !okL {
			return nil, 0, 0, false
		}
		caching.CacheStruct(file, v)
	}
	return projects, limit, offset, true
}

//GetProjects get all the projects with authentication auth
func GetProjects(w http.ResponseWriter, r *http.Request, auth string) []Project {
	var projects []Project
	projectFileName := auth + globals.PROJIDFILE
	//First see if project already exist
	status, filepro := caching.ShouldFileCache(projectFileName, globals.PROJIDDIR)
	defer filepro.Close()
	if status == globals.Error || status == globals.DirFail {
		http.Error(w, "Failed to create a file", http.StatusInternalServerError)
		return nil
	}
	if status == globals.Exist {
		//The file exist
		//We read from file
		err := caching.ReadFile(filepro, &projects)
		if err != nil {
			errmsg := "The Failed Reading from file with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return nil
		}
	} else {
		//Else we need to query to get it
		for i := 0; i < globals.MAXPAGE; i++ {
			var subProj []Project
			query := globals.GITAPI + globals.PROJQ + globals.PAGEQ + strconv.Itoa(i+1)
			err := apiGetCall(w, query, auth, &subProj)
			if err != nil {
				//The API call has failed
				return nil
			}
			//When it's empty we no longer need to do calls
			if len(subProj) == 0 {
				break
			}
			projects = append(projects, subProj...)
		}
		caching.CacheStruct(filepro, projects)

	}
	return projects
}
func findIssuesForProject(project PayloadIssue, auth string,
	w http.ResponseWriter, r *http.Request) ([]Issue, error) {
	//Get projects
	projects := GetProjects(w, r, auth)
	//Find project ID
	projID := "nil"
	for i := range projects {
		if project.ProjectName == projects[i].NamePath {
			projID = strconv.Itoa(projects[i].ID)
			break
		}
	}
	if projID == "nil" {
		return nil, errors.New("No project provided or found")
	}

	return findIssues(projID, auth, w, r), nil
}
func findLabelsInIssues(issues []Issue, auth bool) Labels {
	var labels []Label
	//String map to find duplicates labels
	dupFreq := make(map[string]int)
	for _, issue := range issues {
		for _, label := range issue.Labels {
			dupFreq[label]++
			//first time
			if dupFreq[label] == 1 {
				labels = append(labels, Label{label, 1})
			}
		}
	}
	//Give the the frequency
	for i, label := range labels {
		labels[i].Count = dupFreq[label.Label]
	}
	//Sort it based on frequency
	sort.SliceStable(labels, func(i, j int) bool {
		return dupFreq[labels[i].Label] > dupFreq[labels[j].Label]
	})
	return Labels{labels, auth}
}
func findAuthorsInIssues(issues []Issue, auth bool) Users {
	var users []User
	//String map to find duplicates labels
	dupFreq := make(map[string]int)
	for _, issue := range issues {
		name := issue.Author.Username
		dupFreq[name]++
		//first occurrence
		if dupFreq[name] == 1 {
			users = append(users, User{name, 1})
		}
	}
	//Give the the frequency
	for i := range users {
		users[i].Count = dupFreq[users[i].Username]
	}
	//Sort it based on frequency
	sort.SliceStable(users, func(i, j int) bool {
		return dupFreq[users[i].Username] > dupFreq[users[j].Username]
	})
	return Users{users, auth}
}

//EventOK Checks if the webhooks are valid
func EventOK(event string) bool {
	switch event {
	case string(globals.CommitE):
		return true
	case string(globals.LanguagesE):
		return true
	case string(globals.IssuesE):
		return true
	case string(globals.StatusE):
		return true
	default:
		return false
	}
}

//GetPayload get a payload from body
//Returns error if failed, and a bool that represent if v got filled or not
func GetPayload(r *http.Request, v interface{}) (bool, error) {
	fmt.Println("THE BODY: ", r.Body)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, err
	}
	if len(body) > 0 {
		err = json.Unmarshal(body, &v)
		if err != nil {
			return false, err
		}
		return true, err
	}
	return false, nil
}
func findIssues(projID string, auth string,
	w http.ResponseWriter, r *http.Request) []Issue {
	var issues []Issue
	for i := 0; i < globals.MAXPAGE; i++ {
		var subissues []Issue
		query := globals.GITAPI + globals.PROJQ + projID + "/issues" + globals.PAGEQ + strconv.Itoa(i+1)
		err := apiGetCall(w, query, auth, &subissues)
		if err != nil {

		}
		if len(subissues) == 0 {
			break
		}
		fmt.Println("Issue:", subissues[0].Author)
		issues = append(issues, subissues...)
	}
	return issues
}
