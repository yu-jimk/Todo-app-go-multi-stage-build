package handler

// Create (POST)
type createTodoRequest struct {
	Title string `json:"title"`
}

// UpdateTitle (PATCH)
type updateTitleRequest struct {
	Title string `json:"title"`
}

// UpdateCompleted (PATCH)
type updateCompletedRequest struct {
	Completed bool `json:"completed"`
}
