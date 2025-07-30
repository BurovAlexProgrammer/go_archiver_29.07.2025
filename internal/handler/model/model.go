package model

type TaskStatusResponse struct {
	Status string
}

type ErrorResponse struct {
	Code int
}

type AddFileRequest struct {
	URL string `json:"url"`
}
