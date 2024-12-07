package domain

import "errors"

// Common errors
var (
	ErrNotFound           = errors.New("resource not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidInput       = errors.New("invalid input")
	ErrConflict           = errors.New("resource already exists")
)

// User specific errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already registered")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidEmail    = errors.New("invalid email format")
)

// Attendance specific errors
var (
	ErrAttendanceNotFound      = errors.New("attendance not found")
	ErrAttendanceAlreadyMarked = errors.New("attendance already marked for today")
	ErrInvalidAttendanceStatus = errors.New("invalid attendance status")
)

// Database specific errors
var (
	ErrDatabase = errors.New("database error")
)
