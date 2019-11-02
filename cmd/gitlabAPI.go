package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"RESTGvkGitLab/api"
	"RESTGvkGitLab/caching"
	"RESTGvkGitLab/firedb"
	"RESTGvkGitLab/globals"
)

func main() {
	fmt.Println("Setting up a cleanup interval")
	//Deletes files that are older than 24 hours every 72 hours
	caching.CleanUpInterval(globals.DeleteInterval, globals.DeleteAge)
	fmt.Println("Starting application:")
	firedb.InitDataBase()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := api.SetupHandlers()
	log.Fatal(http.ListenAndServe(":"+port, mux))
	defer firedb.Client.Close()
}
