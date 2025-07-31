// @title Go Zipper API
// @version 1.0
// @description Сервис архивации файлов по ссылкам
// @host localhost:8080
// @BasePath /
package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go_zipper/docs" // автогенерируемая документация
	"go_zipper/internal/configLoader"
	"go_zipper/internal/handler"
	"go_zipper/internal/repository"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := configLoader.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	taskRepository := repository.NewInMemoryTaskRepository()
	taskHandler := handler.NewTaskHandler(taskRepository)
	router := newRouter(taskHandler)
	startServer(cfg, router)

}

func startServer(cfg *configLoader.AppConfig, router *gin.Engine) {
	srv := http.Server{
		Addr:         cfg.HttpSrv.Address,
		IdleTimeout:  cfg.HttpSrv.Timeout,
		ReadTimeout:  cfg.HttpSrv.Timeout,
		WriteTimeout: cfg.HttpSrv.Timeout,
		Handler:      router,
	}

	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}

func newRouter(h *handler.TaskHandler) *gin.Engine {
	ginMode := os.Getenv(gin.EnvGinMode)
	gin.SetMode(ginMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/tasks/create", h.CreateTask)
	router.POST("/tasks/:id/addFile", h.AddFileToTask)
	router.GET("/tasks/:id/status", h.GetStatus)
	router.GET("/tasks/:id/downloadZip", h.DownloadZip)
	return router
}
