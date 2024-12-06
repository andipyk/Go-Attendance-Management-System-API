package domain

import (
	"context"
	"time"
)

type Attendance struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"` // e.g., "present", "absent", "late"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AttendanceRepository interface {
	Create(ctx context.Context, attendance *Attendance) error
	GetByDate(ctx context.Context, date time.Time) ([]Attendance, error)
	GetByUserID(ctx context.Context, userID string) ([]Attendance, error)
	GetByUserIDAndDate(ctx context.Context, userID string, date time.Time) (*Attendance, error)
	Update(ctx context.Context, attendance *Attendance) error
}

type AttendanceUsecase interface {
	MarkAttendance(ctx context.Context, attendance *Attendance) error
	GetAttendanceByDate(ctx context.Context, date time.Time) ([]Attendance, error)
	GetUserAttendance(ctx context.Context, userID string) ([]Attendance, error)
}
