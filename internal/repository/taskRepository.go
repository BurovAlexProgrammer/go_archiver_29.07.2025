package repository

import (
	"go_zipper/internal/domain/task"
)

type TaskRepository interface {
	CreateTask(task *task.Task)
	GetTask(id int) (*task.Task, error)
	UpdateTask(task *task.Task)
	ActiveTaskCount() int
}
