package api

import (
	"time"
)

//The time used to find the uptime of this program
var startTime time.Time

//Version of this current API
const Version string = "v1"

//GITAPI Root API Call
const GITAPI string = "https://git.gvk.idi.ntnu.no/api/v4/"

//GITREPO used for repo calls
const GITREPO string = "/repository/commits/"

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
