package mailLog

import "time"

type SendMailInput struct {
    To      string `json:"to" binding:"required"`
    Subject string `json:"subject" binding:"required"`
    Body    string `json:"body" binding:"required"`
}


type EmailLogResponse struct {
    ID        uint      `json:"id"`
    ToEmail   string    `json:"to_email"`
    Subject   string    `json:"subject"`
    Status    string    `json:"status"`
    Error     string    `json:"error,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

type ListEmailLogsResponse struct {
    Logs  []EmailLogResponse `json:"logs"`
    Total int                `json:"total"`
}