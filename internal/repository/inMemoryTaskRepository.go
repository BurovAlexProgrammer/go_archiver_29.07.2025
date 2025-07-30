package repository

import (
	"errors"
	"go_zipper/internal/domain"
)

type InMemoryTaskRepository struct {
	tasks       map[int]*domain.Task
	idLastIndex int
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[int]*domain.Task),
	}
}

func (r *InMemoryTaskRepository) CreateTask(task *domain.Task) {
	r.idLastIndex++
	task.ID = r.idLastIndex
	task.Status = domain.Created
	r.tasks[task.ID] = task

}

func (r *InMemoryTaskRepository) GetTask(id int) (*domain.Task, error) {
	task, exists := r.tasks[id]

	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (r *InMemoryTaskRepository) UpdateTask(task *domain.Task) {
	r.tasks[task.ID] = task
}
