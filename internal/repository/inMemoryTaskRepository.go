package repository

import (
	"archive/zip"
	"bytes"
	"errors"
	"go_zipper/internal/domain/fileInfo"
	"go_zipper/internal/domain/task"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"sync"
)

var _ TaskRepository = (*InMemoryTaskRepository)(nil) //impl check

type InMemoryTaskRepository struct {
	tasks       map[int]*task.Task
	idLastIndex int
	rwMutex     sync.RWMutex
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[int]*task.Task),
	}
}

func (i *InMemoryTaskRepository) CreateTask(t *task.Task) {
	defer i.rwMutex.Unlock()
	i.rwMutex.Lock()
	i.idLastIndex++
	t.ID = i.idLastIndex
	t.Status = task.InProgress
	i.tasks[t.ID] = t
}

func (i *InMemoryTaskRepository) GetTask(id int) (task.Task, error) {
	i.rwMutex.RLock()
	defer i.rwMutex.RUnlock()
	tsk, exists := i.tasks[id]

	if !exists {
		return task.Task{}, errors.New("task not found")
	}
	return *tsk, nil
}

func (i *InMemoryTaskRepository) UpdateTask(task *task.Task) {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	i.tasks[task.ID] = task
}

func (i *InMemoryTaskRepository) ActiveTaskCount() int {
	i.rwMutex.RLock()
	defer i.rwMutex.RUnlock()
	count := 0
	for _, t := range i.tasks {
		if t.Status == task.InProgress {
			count++
		}
	}
	return count
}

func (i *InMemoryTaskRepository) AddFileInfo(taskId int, fInfo fileInfo.FileInfo) {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	i.tasks[taskId].Files = append(i.tasks[taskId].Files, fInfo)
}

func (i *InMemoryTaskRepository) SetTaskStatus(taskId int, status task.Status) {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	i.tasks[taskId].Status = status
}

func (i *InMemoryTaskRepository) ZipTask(taskId int) error {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	t := i.tasks[taskId]
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for i, file := range t.Files {
		if file.Status != fileInfo.Available {
			slog.Warn("TaskHandler.zipTask: skipped downloading file", "url", file.URL)
			continue
		}

		resp, err := http.Get(file.URL)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Files[i].Status = fileInfo.NotLoaded
			slog.Warn("TaskHandler.zipTask: cannot download file", "url", file.URL, "err", err)
			continue
		}
		defer resp.Body.Close()

		filename := path.Base(resp.Request.URL.Path)
		if filename == "/" || filename == "" {
			filename = "file_" + strconv.Itoa(i)
		}

		fw, err := zipWriter.Create(filename)
		if err != nil {
			slog.Warn("TaskHandler.zipTask: cannot create zip entry", "file", filename, "err", err)
			continue
		}

		b := resp.Body
		written, err := io.Copy(fw, b)
		slog.Info("TaskHandler.zipTask: archived " + strconv.Itoa(int(written)) + " bytes")
		if err != nil {
			slog.Warn("TaskHandler.zipTask: error writing to zip", "file", filename, "err", err)
			continue
		}

		t.Files[i].Status = fileInfo.Loaded
	}
	err := zipWriter.Close()
	if err != nil {
		msg := "TaskHandler.zipTask: cannot close zipWriter. " + err.Error()
		slog.Error(msg)
		return errors.New(msg)
	}

	t.ZipData = buf.Bytes()
	t.Status = task.Archived
	return nil
}
