package domain

type Task struct {
	ID      int        `json:"id"`
	Files   []FileInfo `json:"files"`
	Status  TaskStatus `json:"status"`
	ZipData []int64    `json:"zipData"`
}

const (
	Created TaskStatus = iota
	Preparing
	Finished
)

type TaskStatus int

type FileInfo struct {
	URL    string
	Status FileInfoStatus
}

type FileInfoStatus string

const (
	Available        FileInfoStatus = "Available"
	NotAllowedFormat FileInfoStatus = "NotAllowedFormat"
	Loaded           FileInfoStatus = "Loaded"
	NotLoaded        FileInfoStatus = "NotLoaded"
)
