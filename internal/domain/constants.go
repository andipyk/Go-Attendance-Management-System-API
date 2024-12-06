package domain

const (
	// User roles
	RoleAdmin = "admin"
	RoleUser  = "user"

	// Attendance status
	StatusPresent = "present"
	StatusAbsent  = "absent"
	StatusLate    = "late"

	// Validation constants
	MinPasswordLength = 6
	MaxPasswordLength = 100
	MaxNameLength     = 255
	MaxEmailLength    = 255

	// Time formats
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

// ValidAttendanceStatuses contains all valid attendance statuses
var ValidAttendanceStatuses = map[string]bool{
	StatusPresent: true,
	StatusAbsent:  true,
	StatusLate:    true,
}

// ValidUserRoles contains all valid user roles
var ValidUserRoles = map[string]bool{
	RoleAdmin: true,
	RoleUser:  true,
}
