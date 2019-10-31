package globals

import (
	"time"
)

//FileMsg A "go enum" for what the shouldChacheFile returned
type FileMsg int

const (
	//Error an error occurred
	Error FileMsg = 0
	//OldRenew The file did exist but is now recreated due to age
	OldRenew FileMsg = 1
	//Created we created the file
	Created FileMsg = 2
	//Exist the file Exist and is recent
	Exist FileMsg = 3
	//DirFail directory failed to create
	DirFail FileMsg = 4
	//Deleted the file has been deleted
	Deleted FileMsg = 5
)

//MAXPAGE how many pages (calls we are going to go through)
const MAXPAGE int = 100

//MAXPERPAGE how many result per page (max taken from gitlab API docs)
const MAXPERPAGE int = 100

//PAGEQ is the page query (100 is MAXPAGE)
const PAGEQ string = "?per_page=100&page="

//StartTime The time used to find the uptime of this program
var StartTime time.Time

//Version of this current API
const Version string = "v1"

//DLIMIT is the default limit
const DLIMIT int64 = 5

//PUBLIC prefix for public authentication
const PUBLIC string = "public"

//GITAPI Root API Call
const GITAPI string = "https://git.gvk.idi.ntnu.no/api/v4/"

//GITREPO used for repo calls
const GITREPO string = "/repository/commits/"

//LANGQ the query used to find programmig languages
const LANGQ string = "/languages/"

//PROJQ the query used to find projects
const PROJQ string = "projects/"

//APIFILES is a sub folder where all files get stored
const APIFILES string = "API_Files/"

//COMMITFILE is the file name where we store the repos and commits
const COMMITFILE string = "_commits.json"

//COMMITDIR is the directory it is stores the commit results
const COMMITDIR string = APIFILES + "commits"

//PROJIDFILE ProjectIDFIle is a json file of all project IDs
const PROJIDFILE string = "_project.json"

//PROJIDDIR directory where the project files will be stored
const PROJIDDIR string = APIFILES + "projects"

//LANGFILE the file name we store the languages data
const LANGFILE string = "_languages.Json"

//LANGDIR the directory where we store the
const LANGDIR string = APIFILES + "languages"

//EventMsg are an "enum" for Post request event
type EventMsg string

const (
	//CommitE a Commit Event
	CommitE EventMsg = "commits"
	//LanguagesE a Language Event
	LanguagesE EventMsg = "languages"
	//IssuesE a Issue Event
	IssuesE EventMsg = "issues"
	//StatusE a Status Event
	StatusE EventMsg = "status"
)

//Webhook fields
const (
	WebhookF string = "webhooks"
	EventF   string = "Event"
	IDF      string = "ID"
	TimeF    string = "Time"
	URLF     string = "URL"
)
