package template

type CreateTemplateInput struct {
    Name    string `json:"name" binding:"required"`
    Subject string `json:"subject" binding:"required"`
    Body    string `json:"body" binding:"required"`
    Type    string `json:"type" binding:"required,oneof=notification marketing transactional"`
}

type TemplateResponse struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    Subject   string `json:"subject"`
    Body      string `json:"body"`
    Type      string `json:"type"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
}

type ListTemplatesResponse struct {
    Templates []TemplateResponse `json:"templates"`
    Total     int                `json:"total"`
}