package repository

import "go_zipper/internal/domain"

type TaskRepository interface {
	CreateTask(task *domain.Task)
	GetTask(id string) (*domain.Task, error)
	UpdateTask(task *domain.Task)
}
