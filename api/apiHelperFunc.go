package api

import (
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

	"../caching"
	"../firedb"
	"../globals"
	"google.golang.org/api/iterator"
)

func activateWebhook(event globals.EventMsg) error {
	iter := firedb.Client.Collection("webhooks").Where("event", "==", event).Documents(firedb.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		//TODO: Send granted request
		fmt.Println(doc.Data())
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
	//Some APIcall when calling for commits return a 404,
	//However, I don't want to throw taht error due to 99% of them working
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
	limittxt := r.FormValue("limit")
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
	offsettxt := r.FormValue("offset")
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
	//We have the project now we need to find the amout of commits for each
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
	query := globals.GITAPI + "projects/"
	var lang []string

	//String map to find duplicates
	dupFreq := make(map[string]int)

	//We have the project now we need to find programming langugues
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
		//There is only Get method to commits
		http.Error(w, "only get method allowed", http.StatusNotImplemented)
		return false
	}
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Expecting format .../", http.StatusBadRequest)
		return false
	}
	return true
}

//Works for commits and langugues
func genericGetHandler(w http.ResponseWriter, r *http.Request, fileName string, fileDir string,
	v interface{}, auth string) (int64, int64, bool) {
	if !isGetRequest(w, r) {
		return 0, 0, false
	}
	//Find the headers
	limit := findLimit(w, r)
	offset := findOffset(w, r)

	status, file := caching.ShouldFileCache(fileName, fileDir)
	if status == globals.Error || status == globals.DirFail {
		http.Error(w, "Failed to create a file", http.StatusInternalServerError)
		return limit, offset, false
	}
	if status == globals.Exist {
		//The file exist
		err := caching.ReadFile(file, &v)
		if err != nil {
			errmsg := "The Failed Reading from file with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return limit, offset, false
		}
		//We have no file
	} else {
		projectFileName := auth + globals.PROJIDFILE
		var projects []Project
		//First see if project already exist
		status, file = caching.ShouldFileCache(projectFileName, globals.PROJIDDIR)
		if status == globals.Error || status == globals.DirFail {
			http.Error(w, "Failed to create a file", http.StatusInternalServerError)
			return limit, offset, false
		}
		if status == globals.Exist {
			//The file exist
			//We read from file
			err := caching.ReadFile(file, &projects)
			if err != nil {
				errmsg := "The Failed Reading from file with error" + err.Error()
				http.Error(w, errmsg, http.StatusInternalServerError)
				return limit, offset, false
			}
		} else {
			//Else we need to query to get it
			for i := 0; i < globals.MAXPAGE; i++ {
				var subProj []Project
				query := globals.GITAPI + globals.PROJQ + globals.PAGEQ + strconv.Itoa(i+1)
				err := apiGetCall(w, query, auth, &subProj)
				if err != nil {
					//The API call has failed
					return limit, offset, false
				}
				//When it's empty we no longer need to do calls
				if len(subProj) == 0 {
					break
				}
				projects = append(projects, subProj...)
			}
			caching.CacheStruct(projectFileName, globals.PROJIDDIR, projects)

		}
		repo, okR := v.(*Repos)
		if okR {
			repo.Repos = subAPICallsForCommits(projects, auth, w)
			v = repo
		}
		lang, okL := v.(*Lang)
		if okL {
			lang.Language = subAPICallsForLang(projects, auth, w)
			v = lang
		}
		if !okR && !okL {
			return 0, 0, false
		}
		caching.CacheStruct(fileName, fileDir, v)
	}
	return limit, offset, true
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
