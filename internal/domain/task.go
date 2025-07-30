package domain

type Task struct {
	ID     string     `json:"id"`
	Files  []FileInfo `json:"files"`
	Status TaskStatus `json:"status"`
}

type TaskStatus struct {
}

type FileInfo struct {
	Url string
}
