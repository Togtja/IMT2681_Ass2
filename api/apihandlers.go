package api

import (
	"time"
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

//NilHandler throws a Bad Request
func NilHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)
}
func CommitsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			fmt.Println("Length is: ", len(parts))
			http.Error(w, "Expecting format .../", http.StatusBadRequest)
			return
		}
		limittxt := r.FormValue("limit")
		limit, err := strconv.ParseInt(limittxt, 10, 32)
		if err != nil && limittxt != ""{
			errmsg :=  "Invalid limit, error: " + err.Error()
			http.Error(w, errmsg, http.StatusBadRequest)
			return;
		}
		if limit <= 0 {
			//Default limit
			limit = 5
		}
		offsettxt := r.FormValue("offset")
		offset, err := strconv.ParseInt(offsettxt, 10, 32)
		if err != nil && offsettxt != ""{
			errmsg :=  "Invalid offset, error: " + err.Error()
			http.Error(w, errmsg, http.StatusBadRequest)
			return
		}
		auth := r.FormValue("auth") //TODO: Validate this
		//TEMP
		fmt.Println(offset, auth)


		//TODO Cache all projects
		resp,  err := http.Get(GITAPI + "projects")
		if err != nil {
			errmsg :=  "The HTTP request failed with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errmsg :=  "The HTTP request failed with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}
		var projects []Project
		err = json.Unmarshal(data, &projects)
		if err != nil {
			errmsg :=  "The HTTP request failed with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(projects)

		return
	}
	//There is only Get method to commits
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
	return
}
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			fmt.Println("Length is: ", len(parts))
			http.Error(w, "Expecting format .../", http.StatusBadRequest)
			return
		}
		gitlab,	 err := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects")
		if err != nil {
			errmsg :=  "The HTTP request failed with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}

		db, err := http.Get("https://restcountries.eu/rest/v2/")
		if err != nil {
			errmsg :=  "The HTTP request failed with error" + err.Error()
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}

		uptime := uptime()
		uptimeString := fmt.Sprintf("%.0f seconds", uptime.Seconds())
		diag := StatusDiag{gitlab.StatusCode, db.StatusCode, uptimeString, Version}
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