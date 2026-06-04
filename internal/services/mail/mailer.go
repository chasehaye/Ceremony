package mail

import (
    "fmt"
    "Ceremony/internal/config"
    "Ceremony/internal/services/templates"
    "github.com/resend/resend-go/v2"
)

func SendResetEmail(to, resetLink string) error {
    html := templates.ResetEmail(resetLink)
    return sendMail(to, "Reset your password", html)
}

func SendVerificationEmail(to, verifyLink string) error {
    html := templates.VerificationEmail(verifyLink)
    return sendMail(to, "Verify your email", html)
}

func sendMail(to, subject, body string) error {
    client := resend.NewClient(config.App.ResendAPIKey)

    params := &resend.SendEmailRequest{
        From:    config.App.SenderAddress,
        To:      []string{to},
        Subject: subject,
        Html:    body,
    }

    _, err := client.Emails.Send(params)
    if err != nil {
        fmt.Printf("RESEND ERROR: %v\n", err)
        return fmt.Errorf("resend send failed: %w", err)
    }

    return nil
}