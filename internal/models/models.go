package models

import (
    "time"
    "gorm.io/gorm"
)

type PasswordReset struct {
    gorm.Model
    UserID    uint      `gorm:"index"`
    Token     string    `gorm:"uniqueIndex"`
    ExpiresAt time.Time `gorm:"index"`
    Used      bool      `gorm:"default:false"`
    User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type EmailVerification struct {
    gorm.Model
    UserID    uint      `gorm:"index"`
    Token     string    `gorm:"uniqueIndex"`
    ExpiresAt time.Time `gorm:"index"`
    Used      bool      `gorm:"default:false"`
    User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type User struct {
    ID        uint        `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time  

    Name        string    `gorm:"type:varchar(255)"`
    Password    string    `gorm:"not null"`
    Token       string    `gorm:"type:text"`
    Email       string    `gorm:"uniqueIndex;type:varchar(255);not null"`
	IsAdmin     bool      `gorm:"default:false"`
    IsVerified  bool      `gorm:"default:false"`
    IsApproved  bool      `gorm:"default:false"`
    IsBanned        bool  `gorm:"default:false"`
    CanCreate       bool  `gorm:"default:false"`

    Memberships []OrganizationMember `gorm:"foreignKey:UserID"`
}




type OrganizationMember struct {
    ID        uint `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time

    OrganizationID uint   `gorm:"uniqueIndex:idx_org_user"`
    UserID         uint   `gorm:"uniqueIndex:idx_org_user"`
    Role           string `gorm:"type:varchar(50);default:'member'"`

    Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    User         User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
// Role: member, admin. Single-use; consumed on accept.
type OrganizationInvite struct {
    gorm.Model
    OrganizationID uint      `gorm:"index"`
    InvitedByID    uint      `gorm:"index"`
    Email          string    `gorm:"type:varchar(255);index"` // optional; empty = anyone with the link
    Role           string    `gorm:"type:varchar(50);default:'member'"`
    Token          string    `gorm:"uniqueIndex;type:varchar(255);not null"`
    ExpiresAt      time.Time `gorm:"index"`
    Used           bool      `gorm:"default:false"`

    Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    InvitedBy    User         `gorm:"foreignKey:InvitedByID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type Organization struct {
    gorm.Model
    Name    string `gorm:"type:varchar(255);not null"`
    Slug    string `gorm:"uniqueIndex;type:varchar(255);not null"`

    Members   []OrganizationMember `gorm:"foreignKey:OrganizationID"`
    Invites   []OrganizationInvite `gorm:"foreignKey:OrganizationID"`
    Apps      []App                `gorm:"foreignKey:OrganizationID"`
    EmailLogs []EmailLog           `gorm:"foreignKey:OrganizationID"`
    Templates []EmailTemplate      `gorm:"foreignKey:OrganizationID"`
    Domains   []Domain             `gorm:"foreignKey:OrganizationID"`
}




type App struct {
    gorm.Model
    OrganizationID uint   `gorm:"index"`
    DomainID       *uint  `gorm:"index"`
    Name        string `gorm:"type:varchar(255);not null"`
    Slug        string `gorm:"uniqueIndex;type:varchar(255);not null"` 
    Description string `gorm:"type:text"`
    APIKey      string `gorm:"uniqueIndex;type:varchar(255)"`
    IsActive    bool   `gorm:"default:true"`
    Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    EmailLogs []EmailLog `gorm:"foreignKey:AppID"`
    Domain       *Domain      `gorm:"foreignKey:DomainID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Type: notification, marketing, transactional
type EmailTemplate struct {
	gorm.Model
	OrganizationID uint   `gorm:"index"`
	Name           string `gorm:"type:varchar(255);not null"`
	Subject        string `gorm:"type:varchar(255);not null"` // can include variables e.g. "Welcome, {{.Name}}"
	Body           string `gorm:"type:text;not null"`         // HTML with template variables
	Type           string `gorm:"type:varchar(50)"`           // notification, marketing, transactional
	IsActive       bool   `gorm:"default:true"`

	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	EmailLogs    []EmailLog   `gorm:"foreignKey:TemplateID"`
}


type EmailLog struct {
    gorm.Model
    OrganizationID uint    `gorm:"index"`
    AppID          uint    `gorm:"index"`
    DomainID       *uint  `gorm:"index"`

    FromDomain     string `gorm:"type:varchar(255);not null"`
    FromAddress    string `gorm:"type:varchar(255);not null"`

    ToEmail    string  `gorm:"type:varchar(255);not null"`
    Subject    string  `gorm:"type:varchar(255);not null"`
    Body       string  `gorm:"type:text"`
    TemplateID *uint   `gorm:"index"`
    Status     string  `gorm:"type:varchar(50);default:'pending'"`
    Error      string  `gorm:"type:text"`
    Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    App          App          `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    Template     *EmailTemplate `gorm:"foreignKey:TemplateID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Status: pending, verified, failed
type Domain struct {
	gorm.Model
	OrganizationID uint   `gorm:"index"`
	Name           string `gorm:"type:varchar(255);not null"`
	ResendDomainID string `gorm:"type:varchar(255)"`  // ID returned by Resend after creation
	Status         string `gorm:"type:varchar(50);default:'pending'"`
	Region         string `gorm:"type:varchar(50)"`   // us-east-1, eu-west-1, etc.

	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Apps         []App        `gorm:"foreignKey:DomainID"`
}

