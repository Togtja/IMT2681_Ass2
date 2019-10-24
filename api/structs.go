package api

type Repos struct{
	Repos []Repo `json:"repos"`
	Auth bool `json:"auth"`
}
type Repo struct{
	Repository string `json:"repository"`
	Commits int `json:"commits"`
}
type Lang struct{
	Language []string `json:"languages"`
	Auth bool `json:"auth"`
}
type Users struct{
	Users []User `json:"users"`
}
type User struct{
	Username string `json:"username"`
	Count int `json:"count"`
}
type StatusDiag struct{
	Gitlab int `json:"gitlab"`
	Database int `json:"database"`
	Uptime string `json:"uptime"`
	Version string `json:"version"`
}
//The lengs of it's array is the number of commits

//Details regarding the commit
type Commit struct{
	Msg string `json:"message"`
	Author string `json:"author_name"`
}
//A single project
type Project struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Owner []Owner `json:"owner"`

}
//Details about an Owner
type Owner struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	State string `json:"state"`
}
