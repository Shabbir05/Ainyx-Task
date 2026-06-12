package models

import "time"

// -----------------------------------------------------------------------
// Request bodies
// -----------------------------------------------------------------------

// CreateUserRequest is the JSON body accepted by POST /users.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	DOB  string `json:"dob"  validate:"required,dob_date"`
}

// UpdateUserRequest is the JSON body accepted by PUT /users/:id.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	DOB  string `json:"dob"  validate:"required,dob_date"`
}

// -----------------------------------------------------------------------
// Response bodies
// -----------------------------------------------------------------------

// UserResponse is returned for create / update operations (no age).
type UserResponse struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

// UserWithAgeResponse is returned for read operations (includes age).
type UserWithAgeResponse struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}

// ListUsersResponse wraps a paginated slice of users.
type ListUsersResponse struct {
	Data  []UserWithAgeResponse `json:"data"`
	Page  int32                 `json:"page"`
	Limit int32                 `json:"limit"`
	Total int64                 `json:"total"`
}

// ErrorResponse is the standard JSON error envelope.
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

// -----------------------------------------------------------------------
// Domain model (internal)
// -----------------------------------------------------------------------

// User is the internal representation of a user (not tied to DB layer).
type User struct {
	ID   uint32
	Name string
	DOB  time.Time
}
