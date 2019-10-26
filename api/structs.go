package api

//Repos Represent the Commit repos
type Repos struct {
	Repos []Repo `json:"repos"`
	Auth  bool   `json:"auth"`
}

//Repo represent a single repositiry
type Repo struct {
	Repository string `json:"repository"`
	Commits    int    `json:"commits"`
}

//Lang Represent the language structure
type Lang struct {
	Language []string `json:"languages"`
	Auth     bool     `json:"auth"`
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

//Labels represents the type =labels
type Labels struct {
	Labels []Label `json:"labels"`
	Auth   bool    `json:"auth"`
}

//Label represent the labels and count being num of issues
type Label struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

//Status is the diganostic and status tools
type Status struct {
	Gitlab   int    `json:"gitlab"`
	Database int    `json:"database"`
	Uptime   string `json:"uptime"`
	Version  string `json:"version"`
}

//The lengs of it's array is the number of commits

//Commit Details regarding the commit
type Commit struct {
	Msg    string `json:"message"`
	Author string `json:"author_name"`
}

//Project details concerning the project
type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Owner   []Owner  `json:"owner"`
	Commits []Commit `json:"commits"`
}

//Owner details about the owner of a Project
type Owner struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	State    string `json:"state"`
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
