package handler

import (
	"github.com/gin-gonic/gin"
	"go_zipper/internal/domain"
	"go_zipper/internal/repository"
	"log/slog"
	"net/http"
)

type TaskHandler struct {
	taskRepository repository.TaskRepository
}

func NewTaskHandler(taskRepository repository.TaskRepository) *TaskHandler {
	return &TaskHandler{taskRepository: taskRepository}
}

func (h TaskHandler) CreateTask(ctx *gin.Context) {
	h.taskRepository.CreateTask(&domain.Task{})
}

func (h TaskHandler) AddFileToTask(ctx *gin.Context) {
	taskId := ctx.Param("id")
	if taskId == "" {
		slog.Error("TaskHandler.AddFileToTask: param id cannot be empty")
		return
	}

	task, err := h.taskRepository.GetTask(taskId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Cannot find task")
		slog.Error(err.Error())
		return
	}

	var newFile domain.FileInfo
	err = ctx.ShouldBindJSON(&newFile)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.Files = append(task.Files, newFile)
}

func (h TaskHandler) GetStatus(ctx *gin.Context) {

}
