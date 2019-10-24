package globals

import (
	"time"
)

//FileMsg A "go enum" for what the shouldChacheFile returned
type FileMsg int

const (
	//Error an error occured
	Error FileMsg = 0
	//OldRenew The file did exist but is now recreated due to age
	OldRenew FileMsg = 1
	//Created we created the file
	Created FileMsg = 2
	//Exist the file Exist and is recent
	Exist FileMsg = 3
	//DirFail directory failed to create
	DirFail FileMsg = 4
)

//StartTime The time used to find the uptime of this program
var StartTime time.Time

//Version of this current API
const Version string = "v1"

//GITAPI Root API Call
const GITAPI string = "https://git.gvk.idi.ntnu.no/api/v4/"

//GITREPO used for repo calls
const GITREPO string = "/repository/commits/"

//REPOFILE is the file name where we store the repos
const REPOFILE string = "RepoAndCommits.json"

//REPODIR is the directory it is stores
const REPODIR string = "Repo"

//EventMsg are an "enum" for Post request event
type EventMsg string

const (
	//CommitE a Commit Event
	CommitE EventMsg = "commit"
	//LanguagesE a Language Event
	LanguagesE EventMsg = "languages"
	//IssuesE a Issue Event
	IssuesE EventMsg = "issues"
	//StatusE a Status Event
	StatusE EventMsg = "status"
)
