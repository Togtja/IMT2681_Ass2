package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

//NilHandler throws a Bad Request
func NilHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

//CommitsHandler handler to find the amout of commits
func CommitsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			fmt.Println("Length is: ", len(parts))
			http.Error(w, "Expecting format .../", http.StatusBadRequest)
			return
		}
		limit := findLimit(w, r)
		offset := findOffset(w, r)
		auth := r.FormValue("auth") //TODO: Validate this
		fmt.Println(auth)
		//TODO Cache all projects
		var projects []Project
		query := GITAPI + "projects/"
		if !apicall(w, query, &projects) {
			//The API call has failed
			return
		}
		var repos []Repo
		//We have the project now we need to find the amout of commits for each
		//project
		for i := range projects {
			subquery := query + strconv.Itoa(projects[i].ID) + GITREPO

			var commits []Commit
			if !apicall(w, subquery, &commits) {
				//The API call has failed
				return
			}
			fmt.Println("Commit nr for", projects[i].ID, "is:", len(commits))
			projects[i].Commits = commits
			repos = append(repos, Repo{projects[i].Name, len(commits)})
		}
		sort.SliceStable(repos, func(i, j int) bool {
			return repos[i].Commits > repos[j].Commits
		})
		var finalRepo Repos
		finalRepo.Repos = repos[offset : limit+offset]
		finalRepo.Auth = false
		json.NewEncoder(w).Encode(finalRepo)

		return
	}
	//There is only Get method to commits
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
	return
}

//StatusHandler get status code from db/ external api and uptime and version for thid API
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			fmt.Println("Length is: ", len(parts))
			http.Error(w, "Expecting format .../", http.StatusBadRequest)
			return
		}
		gitlab, err := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects")
		if err != nil {
			errmsg := "The HTTP request failed with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}

		db, err := http.Get("https://restcountries.eu/rest/v2/")
		if err != nil {
			errmsg := "The HTTP request failed with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}

		uptime := uptime()
		uptimeString := fmt.Sprintf("%.0f seconds", uptime.Seconds())
		diag := Status{gitlab.StatusCode, db.StatusCode, uptimeString, Version}
		json.NewEncoder(w).Encode(diag)
		return
	}
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
	return
}

//Finds current uptime
func uptime() time.Duration {
	return time.Since(startTime)
}

//Start the timer to figure out current uptime
func init() {
	startTime = time.Now()
}

func apicall(w http.ResponseWriter, getReq string, v interface{}) bool {
	resp, err := http.Get(getReq)
	if err != nil {
		errmsg := "The HTTP request failed with error" + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return false
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errmsg := "The HTTP request failed with error" + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return false
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		errmsg := "The HTTP request failed with error" + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return false
	}
	return true
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
		limit = 5
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
