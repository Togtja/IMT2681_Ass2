package api

//Repos Represent the Commit repos
type Repos struct {
	Repos []Repo `json:"repos"`
	Auth  bool   `json:"auth"`
}

//Repo represent a single repository
type Repo struct {
	Repository string `json:"repository"`
	Commits    int    `json:"commits"`
}

//Lang Represent the language structure
type Lang struct {
	Language []string `json:"languages"`
	Auth     bool     `json:"auth"`
}

//Payload the expected payload for Lang Get Method
type Payload struct {
	ProjectName []string `json:"projects"`
}

//Issue The struct of an issue in issues
type Issue struct {
	Labels []string `json:"labels"`
	Author Author   `json:"author"`
}

//Author struct of a author in issues
type Author struct {
	Username string `json:"username"`
}

//Users represents the type=users
type Users struct {
	Users []User `json:"users"`
	Auth  bool   `json:"auth"`
}

//User Represent the user count being num of issues
type User struct {
	Username string `json:"username"`
	Count    int    `json:"count"`
}

//Labels represents the type=labels
type Labels struct {
	Labels []Label `json:"labels"`
	Auth   bool    `json:"auth"`
}

//Label represent the labels and count being num of issues
type Label struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

//Status is the diagnostic and status tools
type Status struct {
	Gitlab   int    `json:"gitlab"`
	Database int    `json:"database"`
	Uptime   string `json:"uptime"`
	Version  string `json:"version"`
}

//Commit Details regarding the commit
type Commit struct {
	Msg    string `json:"message"`
	Author string `json:"author_name"`
}

//Project details concerning the project
type Project struct {
	ID       int      `json:"id"`
	NamePath string   `json:"path_with_namespace"`
	Name     string   `json:"name"`
	Commits  []Commit `json:"commits"`
}

//Webhook the structure of an webhook in db
type Webhook struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	URL   string `json:"url"`
	Time  string `json:"time"`
}

//WebhookGet the webhook a user get when calling for a GET method
type WebhookGet struct {
	ID    int    `json:"id"`
	Event string `json:"event"`
	Time  string `json:"time"`
}

//Invocation  it invokes the registered webhook
type Invocation struct {
	Event  string   `json:"event"`
	Params []string `json:"params"`
	Time   string   `json:"timestamp"`
}
