package main

import "time"

// ==================== MODERATOR ====================

type Moderator struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status    string     `json:"status"`
	Token     string     `json:"token,omitempty"`
	Moderator *Moderator `json:"moderator,omitempty"`
	Role      string     `json:"role,omitempty"`
	Error     string     `json:"error,omitempty"`
	Msg       string     `json:"msg,omitempty"`
}

// ==================== CITY ====================

type City struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type CityRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

// ==================== SCHOOL ====================

type School struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	EmailDomain string `json:"email_domain"`
	CityID      int    `json:"city_id"`
	CityName    string `json:"city_name"`
}

type SchoolRequest struct {
	Name        string `json:"name"`
	CityID      int    `json:"city_id"`
	EmailDomain string `json:"email_domain"`
}

type SchoolUpdateRequest struct {
	Name        string `json:"name"`
	EmailDomain string `json:"email_domain"`
}

// ==================== POST ====================

type PendingPost struct {
	ID                int       `json:"id"`
	Content           string    `json:"content"`
	CreatorID         int       `json:"creator_id"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	LikesCount        int       `json:"likes_count"`
	CreatorFirstName  string    `json:"creator_first_name"`
	CreatorLastName   string    `json:"creator_last_name"`
	CreatorEmail      string    `json:"creator_email"`
	SchoolName        *string   `json:"school_name"`
	CityName          *string   `json:"city_name"`
}

type ReportedPost struct {
	PendingPost
	ReportCount int `json:"report_count"`
}

type AllPost struct {
	ID                int       `json:"id"`
	Content           string    `json:"content"`
	CreatorID         int       `json:"creator_id"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	LikesCount        int       `json:"likes_count"`
	CreatorFirstName  string    `json:"creator_first_name"`
	CreatorLastName   string    `json:"creator_last_name"`
	CreatorEmail      string    `json:"creator_email"`
	SchoolName        *string   `json:"school_name"`
	CityName          *string   `json:"city_name"`
	Status            string    `json:"status"`
}

type SetStatusRequest struct {
	Status string `json:"status"`
}

// ==================== SPOTTED ====================

type PendingSpotted struct {
	ID                int       `json:"id"`
	Content           string    `json:"content"`
	CreatorID         int       `json:"creator_id"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	LikesCount        int       `json:"likes_count"`
	Visibility        int       `json:"visibility"`
	VisibilityDesc    string    `json:"visibility_desc"`
	Color             string    `json:"color"`
	CreatorFirstName  string    `json:"creator_first_name"`
	CreatorLastName   string    `json:"creator_last_name"`
	CreatorEmail      string    `json:"creator_email"`
	SchoolName        *string   `json:"school_name"`
	CityName          *string   `json:"city_name"`
}

type ReportedSpotted struct {
	PendingSpotted
	ReportCount int `json:"report_count"`
}

type AllSpotted struct {
	ID                int       `json:"id"`
	Content           string    `json:"content"`
	CreatorID         int       `json:"creator_id"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	LikesCount        int       `json:"likes_count"`
	Visibility        int       `json:"visibility"`
	VisibilityDesc    string    `json:"visibility_desc"`
	Color             string    `json:"color"`
	CreatorFirstName  string    `json:"creator_first_name"`
	CreatorLastName   string    `json:"creator_last_name"`
	CreatorEmail      string    `json:"creator_email"`
	SchoolName        *string   `json:"school_name"`
	CityName          *string   `json:"city_name"`
	Status            string    `json:"status"`
}

// ==================== USERS ====================

type User struct {
	ID            int     `json:"id"`
	Email         string  `json:"email"`
	PersonalEmail *string `json:"personal_email"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Role          string  `json:"role"`
	SchoolName    *string `json:"school_name"`
	CityName      *string `json:"city_name"`
}

type SetRoleRequest struct {
	Role string `json:"role"`
}

// ==================== GENERIC RESPONSES ====================

type SuccessResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg,omitempty"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Msg    string `json:"msg"`
}

type DataResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}
