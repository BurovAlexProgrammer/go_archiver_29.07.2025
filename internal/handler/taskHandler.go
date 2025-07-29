package handler

import "github.com/gin-gonic/gin"

type TaskHandler struct {
}

func (h TaskHandler) CreateTask(context *gin.Context) {

}

func (h TaskHandler) AddFileToTask(context *gin.Context) {

}

func (h TaskHandler) GetStatus(context *gin.Context) {

}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}
