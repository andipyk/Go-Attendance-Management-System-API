package usecase

import (
	"context"
	"golang-tes/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewUserUsecase(userRepo domain.UserRepository, jwtSecret string) domain.UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) Register(ctx context.Context, user *domain.User) error {
	// Check if email already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return domain.ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Set user ID and hashed password
	user.ID = uuid.New().String()
	user.Password = string(hashedPassword)
	if user.Role == "" {
		user.Role = domain.RoleUser
	}

	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", domain.ErrInvalidCredentials
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *userUsecase) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, user *domain.User) error {
	existingUser, err := u.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return domain.ErrUserNotFound
	}

	// If password is provided, hash it
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		user.Password = existingUser.Password
	}

	return u.userRepo.Update(ctx, user)
}
