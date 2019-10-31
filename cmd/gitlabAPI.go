package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"RESTGvkGitLab/api"
	"RESTGvkGitLab/caching"
	"RESTGvkGitLab/firedb"
)

func main() {
	fmt.Println("Setting up a cleanup interval")
	//Deletes files that are older than 24 hours every 72 hours
	caching.CleanUpInterval(72, 24)
	fmt.Println("Starting application:")
	firedb.InitDataBase()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", api.NilHandler)
	http.HandleFunc("/repocheck/v1/commits/", api.CommitsHandler)
	http.HandleFunc("/repocheck/v1/languages/", api.LangHandler)
	http.HandleFunc("/repocheck/v1/issues/", api.IssueHandler)
	http.HandleFunc("/repocheck/v1/status/", api.StatusHandler)
	http.HandleFunc("/repocheck/v1/webhooks/", api.WebhookHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
	defer firedb.Client.Close()
}
