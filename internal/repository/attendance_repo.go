package repository

import (
	"context"
	"database/sql"
	"golang-tes/internal/domain"
	"time"
)

type mysqlAttendanceRepository struct {
	db *sql.DB
}

func NewMySQLAttendanceRepository(db *sql.DB) domain.AttendanceRepository {
	return &mysqlAttendanceRepository{db: db}
}

func (r *mysqlAttendanceRepository) Create(ctx context.Context, attendance *domain.Attendance) error {
	query := `INSERT INTO attendances (id, user_id, attendance_date, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	attendance.CreatedAt = now
	attendance.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, query,
		attendance.ID,
		attendance.UserID,
		attendance.Date,
		attendance.Status,
		attendance.CreatedAt,
		attendance.UpdatedAt,
	)
	return err
}

func (r *mysqlAttendanceRepository) GetByDate(ctx context.Context, date time.Time) ([]domain.Attendance, error) {
	query := `SELECT id, user_id, attendance_date, status, created_at, updated_at 
			  FROM attendances 
			  WHERE DATE(attendance_date) = DATE(?)`

	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []domain.Attendance
	for rows.Next() {
		var attendance domain.Attendance
		err := rows.Scan(
			&attendance.ID,
			&attendance.UserID,
			&attendance.Date,
			&attendance.Status,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, attendance)
	}
	return attendances, nil
}

func (r *mysqlAttendanceRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Attendance, error) {
	query := `SELECT id, user_id, attendance_date, status, created_at, updated_at 
			  FROM attendances 
			  WHERE user_id = ?
			  ORDER BY attendance_date DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []domain.Attendance
	for rows.Next() {
		var attendance domain.Attendance
		err := rows.Scan(
			&attendance.ID,
			&attendance.UserID,
			&attendance.Date,
			&attendance.Status,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, attendance)
	}
	return attendances, nil
}

func (r *mysqlAttendanceRepository) GetByUserIDAndDate(ctx context.Context, userID string, date time.Time) (*domain.Attendance, error) {
	query := `SELECT id, user_id, attendance_date, status, created_at, updated_at 
			  FROM attendances 
			  WHERE user_id = ? AND DATE(attendance_date) = DATE(?)`

	attendance := &domain.Attendance{}
	err := r.db.QueryRowContext(ctx, query, userID, date).Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.Date,
		&attendance.Status,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return attendance, nil
}

func (r *mysqlAttendanceRepository) Update(ctx context.Context, attendance *domain.Attendance) error {
	query := `UPDATE attendances 
			  SET status = ?, updated_at = ? 
			  WHERE id = ?`

	attendance.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query,
		attendance.Status,
		attendance.UpdatedAt,
		attendance.ID,
	)
	return err
}
