package handler

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"go_zipper/internal/domain"
	"go_zipper/internal/handler/dto"
	"go_zipper/internal/repository"
	"io"
	"log/slog"
	"net/http"
	"path"
	"slices"
	"strconv"
)

type TaskHandler struct {
	taskRepository repository.TaskRepository
}

func NewTaskHandler(taskRepository repository.TaskRepository) *TaskHandler {
	return &TaskHandler{taskRepository: taskRepository}
}

// CreateTask godoc
// @Failure 400 {object} dto.ErrorResponse
// @Router /tasks/create [post]
func (h TaskHandler) CreateTask(ctx *gin.Context) {

	task := &domain.Task{}
	//err := ctx.ShouldBindJSON(task)
	//if err != nil {
	//	msg := "TaskHandler.CreateTask: Cannot parse task json"
	//	slog.Error(msg)
	//	ctx.String(http.StatusBadRequest, msg)
	//}
	h.taskRepository.CreateTask(task)
	ctx.IndentedJSON(http.StatusOK, task)
}

// AddFileToTask godoc
// @Failure 400 {object} dto.ErrorResponse
// @Router /tasks/{id}/addFile [post]
// @Param id path string true "ID задачи"
// @Param fileInfo body dto.AddFileRequest true "fileInfoJson"
func (h TaskHandler) AddFileToTask(ctx *gin.Context) {
	taskId, err := h.GetIdParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "TaskHandler.AddFileToTask: "+err.Error())
		return
	}

	task, err := h.taskRepository.GetTask(taskId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "TaskHandler.AddFileToTask: cannot find task")
		slog.Error(err.Error())
		return
	}

	var fileReqData dto.AddFileRequest
	err = ctx.ShouldBindJSON(&fileReqData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "TaskHandler.AddFileToTask:" + err.Error()})
		return
	}

	resp, err := http.Head(fileReqData.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "TaskHandler.AddFileToTask: unable to access file URL"})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	allowedTypes := []string{
		"application/pdf",
		"image/jpeg",
	}

	if slices.Contains(allowedTypes, contentType) == false {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "TaskHandler.AddFileToTask: file type not allowed (only .pdf or .jpeg)"})
		return
	}

	if len(task.Files) >= 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "TaskHandler.AddFileToTask: task already has 3 files"})
		return
	}

	isAlreadyExist := slices.IndexFunc(task.Files, func(f domain.FileInfo) bool { return f.URL == fileReqData.URL }) != -1
	if isAlreadyExist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "TaskHandler.AddFileToTask: file URL already added"})
		return
	}

	newFileInfo := domain.FileInfo{URL: fileReqData.URL}
	task.Files = append(task.Files, newFileInfo)
	ctx.JSON(http.StatusOK, gin.H{"message": "TaskHandler.AddFileToTask: file URL [" + newFileInfo.URL + "] added successfully"})
}

// GetStatus godoc
// @Success 200 {object} dto.TaskStatusResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /tasks/{id}/status [get]
// @Param id path string true "ID задачи"
func (h TaskHandler) GetStatus(ctx *gin.Context) {
	taskId, err := h.GetIdParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "TaskHandler.AddFileToTask: "+err.Error())
		return
	}

	task, err := h.taskRepository.GetTask(taskId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "TaskHandler.AddFileToTask: cannot find task")
		slog.Error(err.Error())
		return
	}

	if len(task.Files) >= 2 {
		h.zipTask(ctx, task)
	}

	ctx.IndentedJSON(http.StatusOK, task)
}

func (h TaskHandler) zipTask(ctx *gin.Context, task *domain.Task) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for i, file := range task.Files {
		resp, err := http.Get(file.URL)
		if err != nil || resp.StatusCode != http.StatusOK {
			task.Files[i].Status = domain.NotLoaded
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
	}
	err := zipWriter.Close()
	if err != nil {
		msg := "TaskHandler.zipTask: cannot close zipWriter. " + err.Error()
		slog.Error(msg)
		ctx.JSON(http.StatusBadRequest, msg)
		return
	}

	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", "attachment; filename=zip_task_"+strconv.Itoa(task.ID)+".zip")
	ctx.Data(http.StatusOK, "application/zip", buf.Bytes())
}

func (h TaskHandler) GetIdParam(ctx *gin.Context) (int, error) {
	taskIdStr := ctx.Param("id")
	if taskIdStr == "" {
		msg := "TaskHandler.GetIdParam: param id cannot be empty"
		slog.Error(msg)
		return -1, errors.New(msg)
	}

	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		msg2 := "TaskHandler.GetIdParam: cannot parse id"
		slog.Error(msg2)
		return -1, errors.New(msg2)
	}

	return taskId, nil
}
