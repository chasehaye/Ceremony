package mailLog

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "Ceremony/internal/dtos"
    "Ceremony/internal/models"
    "Ceremony/internal/services/mail"
)

func Send(c *gin.Context, db *gorm.DB) {
    appID := c.MustGet("appID").(uint)
    // userID := c.MustGet("userID").(uint)

    var input SendMailInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{
            Error: "Invalid input",
            Details: map[string]string{
                "to":      "Recipient email is required",
                "subject": "Subject is required",
                "body":    "Body is required",
            },
        })
        return
    }

    log := models.EmailLog{
        AppID:   appID,
        // UserID:  userID,
        ToEmail: input.To,
        Subject: input.Subject,
        Body:    input.Body,
        Status:  "pending",
    }
    db.Create(&log)

    if err := mail.SendMail(input.To, input.Subject, input.Body); err != nil {
        db.Model(&log).Updates(map[string]interface{}{
            "status": "failed",
            "error":  err.Error(),
        })
        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to send email"})
        return
    }

    db.Model(&log).Update("status", "sent")
    c.JSON(http.StatusOK, gin.H{"message": "Email sent"})
}


func Logs(c *gin.Context, db *gorm.DB) {
    uid := c.MustGet("userID").(uint)

    var logs []models.EmailLog

    if err := db.
        Where("user_id = ?", uid).
        Order("created_at DESC").
        Find(&logs).Error; err != nil {

        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{
            Error: "Failed to fetch logs",
        })
        return
    }

    items := make([]EmailLogResponse, len(logs))

    for i, log := range logs {
        items[i] = EmailLogResponse{
            ID:        log.ID,
            ToEmail:   log.ToEmail,
            Subject:   log.Subject,
            Status:    log.Status,
            Error:     log.Error,
            CreatedAt: log.CreatedAt,
        }
    }

    c.JSON(http.StatusOK, ListEmailLogsResponse{
        Logs:  items,
        Total: len(items),
    })
}