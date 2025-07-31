package domain

import "go_zipper/internal/domain/fileInfo"

type Task struct {
	ID      int                 `json:"id"`
	Files   []fileInfo.FileInfo `json:"files"`
	Status  TaskStatus          `json:"status"`
	ZipData []int64             `json:"zipData"`
}

const (
	Created TaskStatus = iota
	Preparing
	Finished
)

type TaskStatus int
