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