package robot

type SubmitCommandsRequest struct {
	Commands string `json:"commands"`
}

type SubmitCommandsResponse struct {
	TaskID string `json:"taskID"`
}

type CancelTaskResponse struct {
	Message string `json:"message"`
}
