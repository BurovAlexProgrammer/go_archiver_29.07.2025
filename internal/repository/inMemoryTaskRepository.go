package repository

import (
	"errors"
	"go_zipper/internal/domain/task"
)

var _ TaskRepository = (*InMemoryTaskRepository)(nil) //impl check

type InMemoryTaskRepository struct {
	tasks       map[int]*task.Task
	idLastIndex int
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[int]*task.Task),
	}
}

func (r *InMemoryTaskRepository) CreateTask(t *task.Task) {
	r.idLastIndex++
	t.ID = r.idLastIndex
	t.Status = task.InProgress
	r.tasks[t.ID] = t
}

func (r *InMemoryTaskRepository) GetTask(id int) (*task.Task, error) {
	tsk, exists := r.tasks[id]

	if !exists {
		return nil, errors.New("task not found")
	}
	return tsk, nil
}

func (r *InMemoryTaskRepository) UpdateTask(task *task.Task) {
	r.tasks[task.ID] = task
}

func (r *InMemoryTaskRepository) ActiveTaskCount() int {
	count := 0
	for _, t := range r.tasks {
		if t.Status == task.InProgress {
			count++
		}
	}
	return count
}
