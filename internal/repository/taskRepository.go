package repository

import (
	"go_zipper/internal/domain/fileInfo"
	"go_zipper/internal/domain/task"
)

type TaskRepository interface {
	CreateTask(task *task.Task)
	GetTask(id int) (task.Task, error)
	UpdateTask(task *task.Task)
	ActiveTaskCount() int
	AddFileInfo(taskId int, url fileInfo.FileInfo)
	SetTaskStatus(taskId int, status task.Status)
	ZipTask(taskId int) error
}
