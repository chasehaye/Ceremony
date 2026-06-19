package org

import "time"

type CreateOrgInput struct {
	Name string `json:"name" binding:"required"`
}

type OrgResponse struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type OrgWithRoleResponse struct {
	Organization OrgResponse `json:"organization"`
	Role         string      `json:"role"`
}

type UserOrgsResponse struct {
	Organizations []OrgWithRoleResponse `json:"organizations"`
}

type CreateOrgResponse struct {
	Organization OrgResponse `json:"organization"`
}

type DeleteOrgResponse struct {
	Message string `json:"message"`
}

type RecentLogResponse struct {
    ToEmail   string    `json:"to_email"`
    Subject   string    `json:"subject"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

type StatsResponse struct {
    TotalEmails    int64               `json:"total_emails"`
    SentEmails     int64               `json:"sent_emails"`
    FailedEmails   int64               `json:"failed_emails"`
    PendingEmails  int64               `json:"pending_emails"`
    TotalTemplates int64               `json:"total_templates"`
    ActiveApps     int64               `json:"active_apps"`
    RecentLogs     []RecentLogResponse `json:"recent_logs"`
}






type MemberResponse struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type ListMembersResponse struct {
	Members []MemberResponse `json:"members"`
}

type ChangeRoleInput struct {
	Role string `json:"role"` // "admin" or "member"
}

type CreateInviteInput struct {
	Email string `json:"email" binding:"omitempty,email"` // optional; validated only if present
	Role  string `json:"role"`                            // "member" or "admin"; defaults to member
}

type InviteResponse struct {
	ID        uint      `json:"id,omitempty"`
	Token     string    `json:"token"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
	Link      string    `json:"link"`
}

type InvitePreviewResponse struct {
	OrganizationName string `json:"organization_name"`
	OrganizationSlug string `json:"organization_slug"`
	Role             string `json:"role"`
	Email            string `json:"email,omitempty"`
}

type AcceptInviteResponse struct {
	Message          string `json:"message"`
	OrganizationName string `json:"organization_name"`
	OrganizationSlug string `json:"organization_slug"`
	Role             string `json:"role"`
}

type ListInvitesResponse struct {
	Invites []InviteResponse `json:"invites"`
}
