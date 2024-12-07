package usecase

import (
	"context"
	"golang-tes/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAttendanceRepository is a mock type for domain.AttendanceRepository
type MockAttendanceRepository struct {
	mock.Mock
}

func (m *MockAttendanceRepository) Create(ctx context.Context, attendance *domain.Attendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepository) Update(ctx context.Context, attendance *domain.Attendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepository) GetByUserIDAndDate(ctx context.Context, userID string, date time.Time) (*domain.Attendance, error) {
	args := m.Called(ctx, userID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attendance), args.Error(1)
}

func (m *MockAttendanceRepository) GetByDate(ctx context.Context, date time.Time) ([]domain.Attendance, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]domain.Attendance), args.Error(1)
}

func (m *MockAttendanceRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Attendance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Attendance), args.Error(1)
}

func TestAttendanceUsecase_MarkAttendance(t *testing.T) {
	type testCase struct {
		name          string
		attendance    *domain.Attendance
		mockBehavior  func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, attendance *domain.Attendance)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Success",
			attendance: &domain.Attendance{
				UserID: "test-user-id",
				Status: domain.StatusPresent,
			},
			mockBehavior: func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, attendance *domain.Attendance) {
				today := time.Now().Truncate(24 * time.Hour)
				mockUserRepo.On("GetByID", ctx, attendance.UserID).Return(&domain.User{ID: attendance.UserID}, nil)
				mockAttendRepo.On("GetByUserIDAndDate", ctx, attendance.UserID, today).Return(nil, nil)
				mockAttendRepo.On("Create", ctx, mock.AnythingOfType("*domain.Attendance")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Already Marked",
			attendance: &domain.Attendance{
				UserID: "test-user-id",
				Status: domain.StatusPresent,
			},
			mockBehavior: func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, attendance *domain.Attendance) {
				today := time.Now().Truncate(24 * time.Hour)
				mockUserRepo.On("GetByID", ctx, attendance.UserID).Return(&domain.User{ID: attendance.UserID}, nil)
				mockAttendRepo.On("GetByUserIDAndDate", ctx, attendance.UserID, today).Return(&domain.Attendance{}, nil)
			},
			expectedError: domain.ErrAttendanceAlreadyMarked,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockAttendRepo := new(MockAttendanceRepository)
			mockUserRepo := new(MockUserRepository)
			usecase := NewAttendanceUsecase(mockAttendRepo, mockUserRepo)
			ctx := context.Background()

			// Set mock behavior
			tc.mockBehavior(mockAttendRepo, mockUserRepo, ctx, tc.attendance)

			// Execute
			err := usecase.MarkAttendance(ctx, tc.attendance)

			// Assert
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.attendance.ID)
			}
			mockAttendRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAttendanceUsecase_MarkAttendance_Errors(t *testing.T) {
	type testCase struct {
		name          string
		attendance    *domain.Attendance
		mockBehavior  func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, attendance *domain.Attendance)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Database Error on GetByID",
			attendance: &domain.Attendance{
				UserID: "test-user-id",
				Status: domain.StatusPresent,
			},
			mockBehavior: func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, attendance *domain.Attendance) {
				mockUserRepo.On("GetByID", ctx, attendance.UserID).Return(nil, domain.ErrDatabase)
			},
			expectedError: domain.ErrDatabase,
		},
		{
			name: "Database Error on GetByUserIDAndDate",
			attendance: &domain.Attendance{
				UserID: "test-user-id",
				Status: domain.StatusPresent,
			},
			mockBehavior: func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, attendance *domain.Attendance) {
				mockUserRepo.On("GetByID", ctx, attendance.UserID).Return(&domain.User{ID: attendance.UserID}, nil)
				mockAttendRepo.On("GetByUserIDAndDate", ctx, attendance.UserID, mock.AnythingOfType("time.Time")).Return(nil, domain.ErrDatabase)
			},
			expectedError: domain.ErrDatabase,
		},
		{
			name: "Database Error on Create",
			attendance: &domain.Attendance{
				UserID: "test-user-id",
				Status: domain.StatusPresent,
			},
			mockBehavior: func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, attendance *domain.Attendance) {
				mockUserRepo.On("GetByID", ctx, attendance.UserID).Return(&domain.User{ID: attendance.UserID}, nil)
				mockAttendRepo.On("GetByUserIDAndDate", ctx, attendance.UserID, mock.AnythingOfType("time.Time")).Return(nil, nil)
				mockAttendRepo.On("Create", ctx, mock.AnythingOfType("*domain.Attendance")).Return(domain.ErrDatabase)
			},
			expectedError: domain.ErrDatabase,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockAttendRepo := new(MockAttendanceRepository)
			mockUserRepo := new(MockUserRepository)
			usecase := NewAttendanceUsecase(mockAttendRepo, mockUserRepo)
			ctx := context.Background()

			// Set mock behavior
			tc.mockBehavior(mockAttendRepo, mockUserRepo, ctx, tc.attendance)

			// Execute
			err := usecase.MarkAttendance(ctx, tc.attendance)

			// Assert
			assert.ErrorIs(t, err, tc.expectedError)
			mockAttendRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAttendanceUsecase_GetAttendanceByDate(t *testing.T) {
	type testCase struct {
		name             string
		date             time.Time
		mockAttendances  []domain.Attendance
		expectedError    error
		expectedResponse []domain.Attendance
	}

	tests := []testCase{
		{
			name: "Success",
			date: time.Now(),
			mockAttendances: []domain.Attendance{
				{ID: "1", UserID: "user1"},
				{ID: "2", UserID: "user2"},
			},
			expectedError:    nil,
			expectedResponse: []domain.Attendance{{ID: "1", UserID: "user1"}, {ID: "2", UserID: "user2"}},
		},
		{
			name:            "No Records Found",
			date:            time.Now(),
			mockAttendances: []domain.Attendance{},
			expectedError:   domain.ErrAttendanceNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAttendanceRepo := new(MockAttendanceRepository)
			mockUserRepo := new(MockUserRepository)
			usecase := NewAttendanceUsecase(mockAttendanceRepo, mockUserRepo)
			ctx := context.Background()

			mockAttendanceRepo.On("GetByDate", ctx, tc.date).Return(tc.mockAttendances, nil)

			attendances, err := usecase.GetAttendanceByDate(ctx, tc.date)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, attendances)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, attendances)
			}
			mockAttendanceRepo.AssertExpectations(t)
		})
	}
}

func TestAttendanceUsecase_GetAttendanceByDate_Errors(t *testing.T) {
	type testCase struct {
		name          string
		date          time.Time
		mockBehavior  func(mockAttendRepo *MockAttendanceRepository, ctx context.Context, date time.Time)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Database Error",
			date: time.Now(),
			mockBehavior: func(mockAttendRepo *MockAttendanceRepository, ctx context.Context, date time.Time) {
				mockAttendRepo.On("GetByDate", ctx, date).Return([]domain.Attendance{}, domain.ErrDatabase)
			},
			expectedError: domain.ErrDatabase,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAttendRepo := new(MockAttendanceRepository)
			mockUserRepo := new(MockUserRepository)
			usecase := NewAttendanceUsecase(mockAttendRepo, mockUserRepo)
			ctx := context.Background()

			tc.mockBehavior(mockAttendRepo, ctx, tc.date)

			attendances, err := usecase.GetAttendanceByDate(ctx, tc.date)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Nil(t, attendances)
			mockAttendRepo.AssertExpectations(t)
		})
	}
}

func TestAttendanceUsecase_GetUserAttendance(t *testing.T) {
	type testCase struct {
		name             string
		userID           string
		mockUser         *domain.User
		mockAttendances  []domain.Attendance
		expectedError    error
		expectedResponse []domain.Attendance
	}

	tests := []testCase{
		{
			name:   "Success",
			userID: "test-user-id",
			mockUser: &domain.User{
				ID: "test-user-id",
			},
			mockAttendances: []domain.Attendance{
				{ID: "1", UserID: "test-user-id"},
				{ID: "2", UserID: "test-user-id"},
			},
			expectedError:    nil,
			expectedResponse: []domain.Attendance{{ID: "1", UserID: "test-user-id"}, {ID: "2", UserID: "test-user-id"}},
		},
		{
			name:            "User Not Found",
			userID:          "non-existent-id",
			mockUser:        nil,
			mockAttendances: nil,
			expectedError:   domain.ErrUserNotFound,
		},
		{
			name:   "No Attendance Records",
			userID: "test-user-id",
			mockUser: &domain.User{
				ID: "test-user-id",
			},
			mockAttendances: []domain.Attendance{},
			expectedError:   domain.ErrAttendanceNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAttendanceRepo := new(MockAttendanceRepository)
			mockUserRepo := new(MockUserRepository)
			usecase := NewAttendanceUsecase(mockAttendanceRepo, mockUserRepo)
			ctx := context.Background()

			mockUserRepo.On("GetByID", ctx, tc.userID).Return(tc.mockUser, nil)
			if tc.mockUser != nil {
				mockAttendanceRepo.On("GetByUserID", ctx, tc.userID).Return(tc.mockAttendances, nil)
			}

			attendances, err := usecase.GetUserAttendance(ctx, tc.userID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, attendances)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, attendances)
			}
			mockUserRepo.AssertExpectations(t)
			mockAttendanceRepo.AssertExpectations(t)
		})
	}
}

func TestAttendanceUsecase_GetUserAttendance_Errors(t *testing.T) {
	type testCase struct {
		name          string
		userID        string
		mockBehavior  func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, userID string)
		expectedError error
	}

	tests := []testCase{
		{
			name:   "Database Error on GetByID",
			userID: "test-user-id",
			mockBehavior: func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, userID string) {
				mockUserRepo.On("GetByID", ctx, userID).Return(nil, domain.ErrDatabase)
			},
			expectedError: domain.ErrDatabase,
		},
		{
			name:   "Database Error on GetByUserID",
			userID: "test-user-id",
			mockBehavior: func(mockAttendRepo *MockAttendanceRepository, mockUserRepo *MockUserRepository, ctx context.Context, userID string) {
				mockUserRepo.On("GetByID", ctx, userID).Return(&domain.User{ID: userID}, nil)
				mockAttendRepo.On("GetByUserID", ctx, userID).Return([]domain.Attendance{}, domain.ErrDatabase)
			},
			expectedError: domain.ErrDatabase,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAttendRepo := new(MockAttendanceRepository)
			mockUserRepo := new(MockUserRepository)
			usecase := NewAttendanceUsecase(mockAttendRepo, mockUserRepo)
			ctx := context.Background()

			tc.mockBehavior(mockAttendRepo, mockUserRepo, ctx, tc.userID)

			attendances, err := usecase.GetUserAttendance(ctx, tc.userID)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Nil(t, attendances)
			mockAttendRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAttendanceUsecase_MarkAttendance_DefaultStatus(t *testing.T) {
	mockAttendRepo := new(MockAttendanceRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := NewAttendanceUsecase(mockAttendRepo, mockUserRepo)
	ctx := context.Background()

	attendance := &domain.Attendance{
		UserID: "test-user-id",
		// Status is intentionally empty to test default status
	}

	mockUserRepo.On("GetByID", ctx, attendance.UserID).Return(&domain.User{ID: attendance.UserID}, nil)
	mockAttendRepo.On("GetByUserIDAndDate", ctx, attendance.UserID, mock.AnythingOfType("time.Time")).Return(nil, nil)
	mockAttendRepo.On("Create", ctx, mock.AnythingOfType("*domain.Attendance")).Return(nil)

	err := usecase.MarkAttendance(ctx, attendance)
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusPresent, attendance.Status)
	mockAttendRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestAttendanceUsecase_GetAttendanceByDate_DatabaseError(t *testing.T) {
	mockAttendRepo := new(MockAttendanceRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := NewAttendanceUsecase(mockAttendRepo, mockUserRepo)
	ctx := context.Background()
	date := time.Now()

	mockAttendRepo.On("GetByDate", ctx, date).Return([]domain.Attendance{}, domain.ErrDatabase)

	attendances, err := usecase.GetAttendanceByDate(ctx, date)
	assert.ErrorIs(t, err, domain.ErrDatabase)
	assert.Nil(t, attendances)
	mockAttendRepo.AssertExpectations(t)
}

func TestAttendanceUsecase_MarkAttendance_UserNotFound(t *testing.T) {
	mockAttendRepo := new(MockAttendanceRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := NewAttendanceUsecase(mockAttendRepo, mockUserRepo)
	ctx := context.Background()

	attendance := &domain.Attendance{
		UserID: "non-existent-id",
		Status: domain.StatusPresent,
	}

	mockUserRepo.On("GetByID", ctx, attendance.UserID).Return(nil, nil)

	err := usecase.MarkAttendance(ctx, attendance)
	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	mockUserRepo.AssertExpectations(t)
}

func TestAttendanceUsecase_GetUserAttendance_DatabaseErrorOnGetByID(t *testing.T) {
	mockAttendRepo := new(MockAttendanceRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := NewAttendanceUsecase(mockAttendRepo, mockUserRepo)
	ctx := context.Background()

	userID := "test-id"
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, domain.ErrDatabase)

	attendances, err := usecase.GetUserAttendance(ctx, userID)
	assert.ErrorIs(t, err, domain.ErrDatabase)
	assert.Nil(t, attendances)
	mockUserRepo.AssertExpectations(t)
}

func TestAttendanceUsecase_GetUserAttendance_DatabaseErrorOnGetByUserID(t *testing.T) {
	mockAttendRepo := new(MockAttendanceRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := NewAttendanceUsecase(mockAttendRepo, mockUserRepo)
	ctx := context.Background()

	userID := "test-id"
	mockUserRepo.On("GetByID", ctx, userID).Return(&domain.User{ID: userID}, nil)
	mockAttendRepo.On("GetByUserID", ctx, userID).Return([]domain.Attendance{}, domain.ErrDatabase)

	attendances, err := usecase.GetUserAttendance(ctx, userID)
	assert.ErrorIs(t, err, domain.ErrDatabase)
	assert.Nil(t, attendances)
	mockUserRepo.AssertExpectations(t)
	mockAttendRepo.AssertExpectations(t)
}
