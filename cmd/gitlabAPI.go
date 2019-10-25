package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"../api"
	"../firedb"
)

func main() {
	fmt.Println("Starting application:")
	firedb.Test()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", api.NilHandler)
	http.HandleFunc("/repocheck/v1/commits/", api.CommitsHandler)
	http.HandleFunc("/repocheck/v1/languages/", api.LangHandler)
	//http.HandleFunc("/repocheck/v1/issues/", imt2681ass2.CountryHandler)
	http.HandleFunc("/repocheck/v1/status/", api.StatusHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
