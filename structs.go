package imt2681ass2

type Repos struct{
	Repos []Repo `json:"repos"`
	auth bool `json:"auth"`
}
type Repo struct{
	repository string `json:"repository"`
	commits int `json:"commits"`
}
type Lang struct{
	language []string `json:"languages"`
	auth bool `json:"auth"`
}
type Users struct{
	Users []User `json:"users"`
}
type User struct{
	username string `json:"username"`
	count int `json:"count"`
}
type Status struct{
	Gitlab int `json:"gitlab"`
	Database int `json:"database"`
	Uptime string `json:"uptime"`
	Version string `json:"version"`
}