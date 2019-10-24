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

//User Represent the user
type User struct {
	Username string `json:"username"`
	Count    int    `json:"count"`
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
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Owner []Owner `json:"owner"`
}

//Owner details about the owner of a Project
type Owner struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	State    string `json:"state"`
}
