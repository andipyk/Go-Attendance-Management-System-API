package usecase

import (
	"context"
	"golang-tes/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock type for domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestUserUsecase_Register(t *testing.T) {
	type testCase struct {
		name          string
		user          *domain.User
		mockBehavior  func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Success",
			user: &domain.User{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User) {
				mockRepo.On("GetByEmail", ctx, user.Email).Return(nil, nil)
				mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Email Already Exists",
			user: &domain.User{
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User) {
				mockRepo.On("GetByEmail", ctx, user.Email).Return(&domain.User{}, nil)
			},
			expectedError: domain.ErrEmailExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			usecase := NewUserUsecase(mockRepo, "test-secret")
			ctx := context.Background()

			// Set mock behavior
			tc.mockBehavior(mockRepo, ctx, tc.user)

			// Execute
			err := usecase.Register(ctx, tc.user)

			// Assert
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.user.ID)
				assert.NotEqual(t, "password123", tc.user.Password)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_Login(t *testing.T) {
	type testCase struct {
		name          string
		email         string
		password      string
		mockBehavior  func(mockRepo *MockUserRepository, ctx context.Context, email string)
		expectedError error
		expectToken   bool
	}

	tests := []testCase{
		{
			name:     "Success",
			email:    "test@example.com",
			password: "password123",
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, email string) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockRepo.On("GetByEmail", ctx, email).Return(&domain.User{
					ID:       "test-id",
					Email:    email,
					Password: string(hashedPassword),
					Role:     domain.RoleUser,
				}, nil)
			},
			expectedError: nil,
			expectToken:   true,
		},
		{
			name:     "Invalid Credentials",
			email:    "wrong@example.com",
			password: "wrongpass",
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, email string) {
				mockRepo.On("GetByEmail", ctx, email).Return(nil, nil)
			},
			expectedError: domain.ErrInvalidCredentials,
			expectToken:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			usecase := NewUserUsecase(mockRepo, "test-secret")
			ctx := context.Background()

			// Set mock behavior
			tc.mockBehavior(mockRepo, ctx, tc.email)

			// Execute
			token, err := usecase.Login(ctx, tc.email, tc.password)

			// Assert
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_GetProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo, "test-secret")
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := "test-id"
		expectedUser := &domain.User{
			ID:    userID,
			Email: "test@example.com",
			Name:  "Test User",
		}

		mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

		user, err := usecase.GetProfile(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := "non-existent-id"

		mockRepo.On("GetByID", ctx, userID).Return(nil, nil)

		user, err := usecase.GetProfile(ctx, userID)
		assert.Equal(t, domain.ErrUserNotFound, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_UpdateProfile(t *testing.T) {
	type testCase struct {
		name          string
		user          *domain.User
		mockBehavior  func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Success Without Password Change",
			user: &domain.User{
				ID:    "test-id",
				Name:  "Updated Name",
				Email: "test@example.com",
			},
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User) {
				mockRepo.On("GetByID", ctx, user.ID).Return(&domain.User{
					ID:       user.ID,
					Password: "existing-hashed-password",
				}, nil)
				mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Success With Password Change",
			user: &domain.User{
				ID:       "test-id",
				Name:     "Updated Name",
				Email:    "test@example.com",
				Password: "newpassword123",
			},
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User) {
				mockRepo.On("GetByID", ctx, user.ID).Return(&domain.User{
					ID:       user.ID,
					Password: "existing-hashed-password",
				}, nil)
				mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "User Not Found",
			user: &domain.User{
				ID: "non-existent-id",
			},
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User) {
				mockRepo.On("GetByID", ctx, user.ID).Return(nil, nil)
			},
			expectedError: domain.ErrUserNotFound,
		},
		{
			name: "Database Error",
			user: &domain.User{
				ID: "test-id",
			},
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User) {
				mockRepo.On("GetByID", ctx, user.ID).Return(nil, domain.ErrDatabase)
			},
			expectedError: domain.ErrDatabase,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			usecase := NewUserUsecase(mockRepo, "test-secret")
			ctx := context.Background()

			// Set mock behavior
			tc.mockBehavior(mockRepo, ctx, tc.user)

			// Execute
			err := usecase.UpdateProfile(ctx, tc.user)

			// Assert
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// Add error cases for existing tests
func TestUserUsecase_Register_Errors(t *testing.T) {
	tests := []struct {
		name          string
		user          *domain.User
		mockBehavior  func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User)
		expectedError error
	}{
		{
			name: "Database Error on GetByEmail",
			user: &domain.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User) {
				mockRepo.On("GetByEmail", ctx, user.Email).Return(nil, domain.ErrDatabase)
			},
			expectedError: domain.ErrDatabase,
		},
		{
			name: "Database Error on Create",
			user: &domain.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockBehavior: func(mockRepo *MockUserRepository, ctx context.Context, user *domain.User) {
				mockRepo.On("GetByEmail", ctx, user.Email).Return(nil, nil)
				mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(domain.ErrDatabase)
			},
			expectedError: domain.ErrDatabase,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			usecase := NewUserUsecase(mockRepo, "test-secret")
			ctx := context.Background()

			tc.mockBehavior(mockRepo, ctx, tc.user)

			err := usecase.Register(ctx, tc.user)
			assert.ErrorIs(t, err, tc.expectedError)
			mockRepo.AssertExpectations(t)
		})
	}
}
