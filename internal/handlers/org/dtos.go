package org

type CreateOrgInput struct {
	Name string `json:"name" binding:"required"`
}

type OrgResponse struct {
	ID   uint   `json:"id"`
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