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
}

type ListUsersResponse struct {
    Users []UserListItem `json:"users"`
    Total int            `json:"total"`
}