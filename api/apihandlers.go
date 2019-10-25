package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"../globals"

	"../caching"
)

//NilHandler throws a Bad Request
func NilHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

//CommitsHandler handler to find the amout of commits
func CommitsHandler(w http.ResponseWriter, r *http.Request) {
	var repo Repos
	repo.Auth = false
	//Default public
	commitFileName := globals.PUBLIC
	auth := r.FormValue("auth")
	if auth != "" {
		//Make a personal json file for authorized users
		//Should be deleted after XX hours/Days
		repo.Auth = true
		commitFileName = auth
	} else {
		auth = globals.PUBLIC
	}
	commitFileName = commitFileName + globals.COMMITFILE
	limit, offset, ok := genericHandler(w, r, commitFileName, globals.COMMITDIR, &repo, auth)
	if ok == false {
		return
	}
	if limit+offset > int64(len(repo.Repos)) {
		if offset >= int64(len(repo.Repos)) {
			offset = 0
		}
		limit = int64(len(repo.Repos)) - offset

	}
	repo.Repos = repo.Repos[offset : limit+offset]
	json.NewEncoder(w).Encode(repo)

	return

}

//LangHandler handles the Programming Language requests
func LangHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only get method allowed", http.StatusNotImplemented)
		return
	}
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Expecting format .../", http.StatusBadRequest)
		return
	}
	//Find the headers
	//limit := findLimit(w, r)
	//offset := findOffset(w, r)

	var lang Lang
	lang.Auth = false
	langFileName := globals.PUBLIC
	auth := r.FormValue("auth")
	if auth != "" {
		//Make a personal json file for authorized users
		//TODO: Should be deleted after XX hours/Days
		lang.Auth = true
		langFileName = auth
	} else {
		auth = globals.PUBLIC
	}
	langFileName = langFileName + globals.LANGFILE
	limit, offset, ok := genericHandler(w, r, langFileName, globals.LANGDIR, &lang, auth)
	if ok == false {
		return
	}
	if limit+offset > int64(len(lang.Language)) {
		if offset >= int64(len(lang.Language)) {
			offset = 0
		}
		limit = int64(len(lang.Language)) - offset

	}
	lang.Language = lang.Language[offset : limit+offset]
	json.NewEncoder(w).Encode(lang)
}

//IssueHandler handles issue request
func IssueHandler(w http.ResponseWriter, r *http.Request) {
	if !isCorrectRequest(w, r) {
		return
	}
	_type := r.FormValue("type")
	if _type == "users" {

	} else if _type == "labels" {

	} else {

	}
}

//StatusHandler get status code from db/ external api and uptime and version for thid API
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only get method allowed", http.StatusNotImplemented)
		return
	}
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Expecting format .../", http.StatusBadRequest)
		return
	}
	gitlab, err := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects")
	if err != nil {
		errmsg := "The HTTP request failed with error" + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	//TODO: FIX THIS (I.E call database)
	db, err := http.Get("https://restcountries.eu/rest/v2/")
	if err != nil {
		errmsg := "The HTTP request failed with error" + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	uptimeString := fmt.Sprintf("%.0f seconds", uptime().Seconds())
	diag := Status{gitlab.StatusCode, db.StatusCode, uptimeString, globals.Version}
	json.NewEncoder(w).Encode(diag)
	return

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
	if auth != globals.PUBLIC {
		request.Header.Set("Private-Token", auth)
	}
	resp, err := client.Do(request)
	if err != nil {
		errmsg := "The HTTP request failed with error: " + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return err
	}
	if resp.StatusCode == 404 {
		return nil
	}
	if resp.StatusCode == 401 {
		//Invalid auth
		http.Error(w, "Invalid Authentication", http.StatusUnauthorized)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errmsg := "The Read of the response failed with error: " + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return err
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		fmt.Println(request, "\n", v, resp)
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
			repos = append(repos, Repo{projects[i].Name, len(projects[i].Commits)})
			m.Unlock()
			wg.Done()
		}(i)
	}
	fmt.Println(&wg)
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

	//We have the project now we need to find the amout of commits for each
	//project
	var wg sync.WaitGroup
	var m = &sync.Mutex{}
	for i := range projects {
		wg.Add(1)
		//Do calls in multithreading
		go func(i int) {
			subquery := query + strconv.Itoa(projects[i].ID) + globals.LANGQ
			fmt.Println(subquery)
			var v interface{}
			err := apiGetCall(w, subquery, auth, &v)
			if err != nil {
				//The API call has failed
				wg.Done()
				return
			}
			data := v.(map[string]interface{})
			var higest float64 = 0
			var language string = ""
			for k, v := range data {
				switch v := v.(type) {
				case float64:
					if v > higest {
						higest = v
						language = k
					}
				default:
					continue
				}
			}
			if language == "" {
				wg.Done()
				return
			}
			//Make sure we don't append at the same time
			m.Lock()
			dupFreq[language]++
			//If it is the first time we seen it
			if dupFreq[language] == 1 {

				lang = append(lang, language)
			}
			m.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()

	sort.SliceStable(lang, func(i, j int) bool {
		return dupFreq[lang[i]] > dupFreq[lang[j]]
	})

	return lang
}
func isCorrectRequest(w http.ResponseWriter, r *http.Request) bool {
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

//Work for commits and langugues
func genericHandler(w http.ResponseWriter, r *http.Request, fileName string, fileDir string,
	v interface{}, auth string) (int64, int64, bool) {
	if !isCorrectRequest(w, r) {
		return 0, 0, false
	}
	//Find the headers
	limit := findLimit(w, r)
	offset := findOffset(w, r)

	file := caching.FileExist(fileName, fileDir)
	if file != nil {
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
		file := caching.FileExist(fileName, fileDir)
		if file != nil {
			//The file exist
			//We read from file
			err := caching.ReadFile(file, &projects)
			if err != nil {
				errmsg := "The Failed Reading from file with error" + err.Error()
				http.Error(w, errmsg, http.StatusInternalServerError)
				return limit, offset, false
			}

		} else {
			fmt.Println("We get here")
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
					fmt.Println("subProj is", len(subProj), "with query", query)
					break
				}
				projects = append(projects, subProj...)
			}
			//Else we need to quary to get it

			fmt.Println("But not here?")
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
