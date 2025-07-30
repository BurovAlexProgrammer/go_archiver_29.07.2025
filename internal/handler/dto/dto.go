package dto

type TaskStatusResponse struct {
	Status string `json:"status" default:"Created"`
}

type ErrorResponse struct {
	Code int `json:"code"`
}

type AddFileRequest struct {
	URL string `json:"url"`
}

type CreateTaskRequest struct {
	//Status string `json:"status" default:"Created"`
}
