package task

import "go_zipper/internal/domain/fileInfo"

type Task struct {
	ID      int                 `json:"id"`
	Files   []fileInfo.FileInfo `json:"files"`
	Status  Status              `json:"status"`
	ZipData []byte              `json:"zipData"`
}

const (
	Undefined Status = iota
	InProgress
	Archived
)

var _ = Undefined //Default value

type Status int
