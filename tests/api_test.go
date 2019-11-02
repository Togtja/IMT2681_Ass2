package tests

import (
	"RESTGvkGitLab/api"
	"RESTGvkGitLab/globals"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

//TestCommitHandler tests the commit handler
func TestCommitHandler(t *testing.T) {
	/*This will fail unless there is test data inside of the tests folder*/
	/***********
	Struct sample
	tests --> API_Files --> commits --> public_commits.json
						--> projects --> public_projects.json
	************/
	//Make sure we ignore age of files so the test files does not get deleted
	globals.DeleteAge = -1
	rr := checkRequest(t, "GET", "/repocheck/v1/commits/", http.StatusOK)
	if rr == nil {
		t.Errorf("Error Expecting  a body")
	}
	repo1 := api.Repo{"TestRepo1", 1074}
	repo2 := api.Repo{"Test/Repo/2", 69}
	var repos []api.Repo
	repos = append(repos, repo1)
	repos = append(repos, repo2)
	expected := api.Repos{repos, false}

	var a api.Repos
	err := json.NewDecoder(rr.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}
	if !reflect.DeepEqual(a, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			a, expected)
	}
}
func TestCommitHandlerNotImplemented(t *testing.T) {
	onlyGetRequest(t, "/repocheck/v1/commits/")
}

func TestLangHandler(t *testing.T) {
	globals.DeleteAge = -1
	rr := checkRequest(t, "GET", "/repocheck/v1/languages/", http.StatusOK)
	if rr == nil {
		t.Errorf("Error Expecting  a body")
	}
	var a api.Lang
	err := json.NewDecoder(rr.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}

	var expected api.Lang
	var langs []string
	langs = append(langs, "Java")
	langs = append(langs, "Go")
	langs = append(langs, "C++")

	expected.Language = langs
	expected.Auth = false

	if !reflect.DeepEqual(a, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			a, expected)
	}

}

//Test that the correct methods are implemented
func TestLangHandlerNotImplemented(t *testing.T) {
	onlyGetRequest(t, "/repocheck/v1/languages/")
}

//TestMainThings Run firedb rest to check if database works
func TestStatusHandlerGet(t *testing.T) {

	rr := checkRequest(t, "GET", "/repocheck/v1/status/", http.StatusOK)
	var expected api.Status
	expected.Gitlab = 200
	expected.Database = 200
	expected.Version = "v1"
	expected.Uptime = "seconds" //Used for sun string comare
	//Used to give a good error message
	expectedTxt := `{"gitlab":200,"database":200,"uptime":"__IGNORE THIS__","version":"v1"}`
	body := strings.TrimSpace(rr.Body.String())

	var a api.Status
	err := json.NewDecoder(rr.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}
	//We do not care for uptime as that will change depending on
	if a.Gitlab != expected.Gitlab || a.Database != expected.Database || a.Version != expected.Version {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, expectedTxt)
	}
	if !strings.Contains(a.Uptime, expected.Uptime) {
		t.Errorf("handler returned unexpected uptime: got %v want x %v (where is x any integer)",
			a.Uptime, expectedTxt)
	}

}

func TestStatusHandlerNotImplemented(t *testing.T) {
	onlyGetRequest(t, "/repocheck/v1/status/")
}
