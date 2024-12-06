package usecase

import (
	"context"
	"golang-tes/internal/domain"
	"time"

	"github.com/google/uuid"
)

type attendanceUsecase struct {
	attendanceRepo domain.AttendanceRepository
	userRepo       domain.UserRepository
}

func NewAttendanceUsecase(attendanceRepo domain.AttendanceRepository, userRepo domain.UserRepository) domain.AttendanceUsecase {
	return &attendanceUsecase{
		attendanceRepo: attendanceRepo,
		userRepo:       userRepo,
	}
}

func (u *attendanceUsecase) MarkAttendance(ctx context.Context, attendance *domain.Attendance) error {
	user, err := u.userRepo.GetByID(ctx, attendance.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	// Check if attendance already exists for today
	today := time.Now().Truncate(24 * time.Hour)
	existing, err := u.attendanceRepo.GetByUserIDAndDate(ctx, attendance.UserID, today)
	if err != nil {
		return err
	}
	if existing != nil {
		return domain.ErrAttendanceAlreadyMarked
	}

	// Set default status if not provided
	if attendance.Status == "" {
		attendance.Status = domain.StatusPresent
	}

	// Continue with marking attendance
	attendance.ID = uuid.New().String()
	attendance.Date = time.Now()

	return u.attendanceRepo.Create(ctx, attendance)
}

func (u *attendanceUsecase) GetAttendanceByDate(ctx context.Context, date time.Time) ([]domain.Attendance, error) {
	attendances, err := u.attendanceRepo.GetByDate(ctx, date)
	if err != nil {
		return nil, err
	}
	if len(attendances) == 0 {
		return nil, domain.ErrAttendanceNotFound
	}
	return attendances, nil
}

func (u *attendanceUsecase) GetUserAttendance(ctx context.Context, userID string) ([]domain.Attendance, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	attendances, err := u.attendanceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(attendances) == 0 {
		return nil, domain.ErrAttendanceNotFound
	}

	return attendances, nil
}
