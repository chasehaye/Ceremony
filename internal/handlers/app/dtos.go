package app

type CreateAppInput struct {
    Name        string `json:"name" binding:"required"`
    Description string `json:"description"`
}

type AppResponse struct {
    ID          uint   `json:"id"`
    Slug        string `json:"slug"`
    Name        string `json:"name"`
    Description string `json:"description"`
    APIKey      string `json:"api_key"`
    IsActive    bool   `json:"is_active"`
    CreatedAt   string `json:"created_at"`
}

type ListAppsResponse struct {
    Apps  []AppResponse `json:"apps"`
    Total int           `json:"total"`
}