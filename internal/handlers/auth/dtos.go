package auth


type RegisterInput struct {
	Name     string `json:"name" binding:"max=255" example:"User Name"`
	Email    string `json:"email" binding:"required,email,max=255" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8,max=72" example:"SecurePass123!"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email,max=255" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8,max=72" example:"SecurePass123!"`
}

type LogOutResponse struct {
	Message  string `json:"message" example:"success"`
}

type ForgotPasswordInput struct {
	Email    string `json:"email" binding:"required,email,max=255" example:"user@example.com"`
}

type ForgotPasswordResponse struct {
	Message  string `json:"message" example:"Check your inbox for a reset link"`
}

type ResetPasswordInput struct {
	Password string `json:"password" binding:"required,min=8,max=72" example:"SecurePass123!"`
}

type ResetPasswordResponse struct {
	Message  string `json:"message" example:"Password updated successfully"`
}

type RegisterResponse struct {
    Message    string `json:"message"`
    IsAdmin    bool   `json:"is_admin"`
    IsVerified bool   `json:"is_verified"`
    IsApproved bool   `json:"is_approved"`
    IsBanned   bool   `json:"is_banned"`
    CanCreate  bool   `json:"can_create"`
    UserEmail  string `json:"user_email"`
    UserName   string `json:"user_name"`
    IsSuperAdmin bool   `json:"is_super_admin"`
}

type LoginResponse struct {
    Message    string `json:"message"`
    IsAdmin    bool   `json:"is_admin"`
    IsVerified bool   `json:"is_verified"`
    IsApproved bool   `json:"is_approved"`
    IsBanned   bool   `json:"is_banned"`
    CanCreate  bool   `json:"can_create"`
    UserEmail  string `json:"user_email"`
    UserName   string `json:"user_name"`
    IsSuperAdmin bool   `json:"is_super_admin"`
}

type GetMeResponse struct {
    ID           uint   `json:"id"`
    UserName     string `json:"user_name"`
    UserEmail    string `json:"user_email"`
    IsAdmin      bool   `json:"is_admin"`
    IsVerified   bool   `json:"is_verified"`
    IsApproved   bool   `json:"is_approved"`
    IsBanned     bool   `json:"is_banned"`
    CanCreate    bool   `json:"can_create"`
    IsSuperAdmin bool   `json:"is_super_admin"`
}