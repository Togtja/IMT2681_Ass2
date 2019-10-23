package imt2681ass2

import (
	"time"
	"strings"
)

//The time used to find the uptime of this program
var startTime time.Time

//Version of this current API
const Version string = "v1"

type EventMsg string

const (
	//The event is message
	Commit EventMsg = "commit"
	//the event is languges
	Languages EventMsg = "languages"
	//The event is issues
	Issues EventMsg = "issues"
	//The event is status
	Status EventMsg = "status"
)