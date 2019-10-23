package main

import (
	"fmt"
	"imt2681ass2"
	"log"
	"net/http"
	"os"
)
func main() {
	fmt.Println("Starting application:")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", imt2681ass2.NilHandler)
	//http.HandleFunc("/repocheck/v1/commits/", imt2681ass2.CountryHandler)
	//http.HandleFunc("/repocheck/v1/languages/", imt2681ass2.SpeciesHandler)
	//http.HandleFunc("/repocheck/v1/issues/", imt2681ass2.CountryHandler)
	http.HandleFunc("/repocheck/v1/status/", imt2681ass2.StatusHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
