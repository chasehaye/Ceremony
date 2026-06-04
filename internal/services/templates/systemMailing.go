package templates

import "fmt"

var baseStyle = `
    font-family: 'Inter', sans-serif;
    background-color: #ffffff;
    color: #111113;
    padding: 40px 20px;
    max-width: 480px;
    margin: 0 auto;
`

var cardStyle = `
    background-color: #f9f9f9;
    border: 1px solid #e5e5e5;
    border-radius: 16px;
    padding: 32px;
    text-align: center;
`

var buttonStyle = `
    background-color: #3d6b4f;
    color: #ffffff;
    padding: 12px 24px;
    text-decoration: none;
    display: inline-block;
    border-radius: 999px;
    font-size: 14px;
    font-weight: 500;
`

var mutedStyle = `
    color: #888888;
    font-size: 12px;
`

var headingStyle = `
    font-size: 20px;
    font-weight: 500;
    color: #111113;
    margin-bottom: 8px;
`

func ResetEmail(resetLink string) string {
	return fmt.Sprintf(`
		<div style="%s">
			<div style="%s">
				<p style="font-size: 13px; color: #7fbf95; margin-bottom: 8px;">● Ceremony</p>
				<h1 style="font-size: 20px; font-weight: 500; margin-bottom: 8px;">Reset your password</h1>
				<p style="%s margin-bottom: 24px;">
					You requested a password reset. Click the button below to continue.
				</p>
				<a href="%s" style="%s">Reset Password</a>
				<p style="%s margin-top: 24px;">
					If the button doesn't work, copy and paste this link:
				</p>
				<p style="word-break: break-all; font-size: 11px; color: #555;">%s</p>
			</div>
		</div>
	`, baseStyle, cardStyle, mutedStyle, resetLink, buttonStyle, mutedStyle, resetLink)
}

func VerificationEmail(verifyLink string) string {
	return fmt.Sprintf(`
		<div style="%s">
			<div style="%s">
				<p style="font-size: 13px; color: #3d6b4f; margin-bottom: 8px;">● Ceremony</p>
				<h1 style="font-size: 20px; font-weight: 500; margin-bottom: 8px;">Verify your email</h1>
				<p style="%s margin-bottom: 24px;">
					Click the button below to verify your email and continue.
				</p>
				<a href="%s" style="%s">Verify Email</a>
				<p style="%s margin-top: 24px;">
					If the button doesn't work, copy and paste this link:
				</p>
				<p style="word-break: break-all; font-size: 11px; color: #555;">%s</p>
			</div>
		</div>
	`, baseStyle, cardStyle, mutedStyle, verifyLink, buttonStyle, mutedStyle, verifyLink)
}