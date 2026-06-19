package admin

import "time"

type UserListItem struct {
    ID         uint      `json:"id"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    Name       string    `json:"name"`
    Email      string    `json:"email"`
    IsAdmin    bool      `json:"is_admin"`
    IsVerified bool      `json:"is_verified"`
    IsApproved bool      `json:"is_approved"`
    IsBanned   bool      `json:"is_banned"`
    CanCreate  bool      `json:"can_create"`
    IsSuperAdmin bool `json:"is_super_admin"`
}

type ListUsersResponse struct {
    Users []UserListItem `json:"users"`
    Total int            `json:"total"`
}