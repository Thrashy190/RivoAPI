package project

type Project struct {
	ID          string `dynamodbav:"id"`
	Name        string `dynamodbav:"name"`
	Description string `dynamodbav:"description"`
	CreatedAt   int64  `dynamodbav:"created_at"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetProjectListResponse struct {
	Projects []Project `json:"projects"`
}
