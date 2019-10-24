package api
import (
	"time"
)

//The time used to find the uptime of this program
var startTime time.Time

//Version of this current API
const Version string = "v1"

//Root API Call
const GITAPI string = "https://git.gvk.idi.ntnu.no/api/v4/"

type EventMsg string

const (
	//The event is message
	CommitE EventMsg = "commit"
	//the event is languges
	LanguagesE EventMsg = "languages"
	//The event is issues
	IssuesE EventMsg = "issues"
	//The event is status
	StatusE EventMsg = "status"
)