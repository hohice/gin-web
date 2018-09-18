package project

const urltag = "projects"

type File struct {
	Content       string `json:"content"`
	Path          string `json:"path"`
	Branch        string `json:"branch"`
	CommitMessage string `json:"commit_message"`
}

type Project struct {
	Files       []File `json:"files"`
	NamespaceId uint64 `json:"namespace_id"`
	Id          uint64 `json:"id"`
	Visibility  string `json:"visibility"`
	Name        string `json:"name"`
}

type ProjectInfo struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}
