package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

		var finalRepo Repos
		finalRepo.Auth = false
		fileName := "Common"
		auth := r.FormValue("auth") //TODO: is sorta done look at TODO furhter down
		if auth != "" {
			finalRepo.Auth = true
			fileName = auth
		}
		fileName = fileName + globals.REPOFILE
		file := caching.FileExist(fileName, globals.REPODIR)
		if file != nil {
			//The gile exist
			data, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
				return
			}
			json.Unmarshal(data, &finalRepo)
			file.Close()
			//We have no file
		} else {
			//TODO Cache all projects
			var projects []Project
			query := globals.GITAPI + "projects/"
			if !apiGetCall(w, query, auth, &projects) {
				//The API call has failed
				return
			}
			var repos []Repo
			//We have the project now we need to find the amout of commits for each
			//project
			var wg sync.WaitGroup
			var m = &sync.Mutex{}
			for i := range projects {
				wg.Add(1)
				//Do calls in multithreading
				go func(i int) {
					subquery := query + strconv.Itoa(projects[i].ID) + globals.GITREPO

					var commits []Commit
					if !apiGetCall(w, subquery, auth, &commits) {
						//The API call has failed
						return
					}
					projects[i].Commits = commits
					//Make sure we don't append at the same time
					m.Lock()
					repos = append(repos, Repo{projects[i].Name, len(commits)})
					m.Unlock()
					//fmt.Println("We have recived data")
					wg.Done()
				}(i)
			}
			wg.Wait()
			sort.SliceStable(repos, func(i, j int) bool {
				return repos[i].Commits > repos[j].Commits
			})

			finalRepo.Repos = repos
			caching.CacheStruct(fileName, globals.REPODIR, finalRepo)
		}
		//Make sure we don't go to high
		if limit+offset > int64(len(finalRepo.Repos)) {
			if offset >= int64(len(finalRepo.Repos)) {
				offset = 0
			}
			limit = int64(len(finalRepo.Repos)) - offset

		}
		finalRepo.Repos = finalRepo.Repos[offset : limit+offset]
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
		diag := Status{gitlab.StatusCode, db.StatusCode, uptimeString, globals.Version}
		json.NewEncoder(w).Encode(diag)
		return
	}
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
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

func apiGetCall(w http.ResponseWriter, getReq string, auth string, v interface{}) bool {
	client := &http.Client{}
	request, err := http.NewRequest("GET", getReq, nil)

	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Private-Token", auth)
	resp, err := client.Do(request)
	if err != nil {
		errmsg := "The HTTP request failed with error: " + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return false
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errmsg := "The Read of the response failed with error: " + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return false
	}
	//TODO: Figure out why it throws an umashal error when auth is invalid
	//TODO: Figure out why Owner is not accpeted
	err = json.Unmarshal(data, &v)
	if err != nil {
		errmsg := "The Unmarshal failed with error: " + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		fmt.Println(data)
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
