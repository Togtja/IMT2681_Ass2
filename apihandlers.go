package imt2681ass2

import (
	"time"
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
)

//NilHandler throws a Bad Request
func NilHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)
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
		gitlab,	 err := http.Get("http://api.gbif.org/v1/")
		if err != nil {
			http.Error(w, "The HTTP request failed with error", http.StatusInternalServerError)
			fmt.Printf("The HTTP request failed with error %s\n", err)
			return
		}

		db, err := http.Get("https://restcountries.eu/rest/v2/")
		if err != nil {
			http.Error(w, "The HTTP request failed with error", http.StatusInternalServerError)
			fmt.Printf("The HTTP request failed with error %s\n", err)
			return
		}

		uptime := uptime()
		uptimeString := fmt.Sprintf("%.0f seconds", uptime.Seconds())
		diag := Status{gitlab.StatusCode, db.StatusCode, uptimeString, Version}
		json.NewEncoder(w).Encode(diag)
		fmt.Println("Sucsess??")
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