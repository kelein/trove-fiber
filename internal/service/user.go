package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	v1 "github.com/kelein/trove-fiber/internal/api/v1"
	"github.com/kelein/trove-fiber/internal/model"
	"github.com/kelein/trove-fiber/internal/repository"
)

// User Service Errors
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserYetExist  = errors.New("user email already exists")
	ErrDatabaseQuery = errors.New("database error when querying")
)

// UserService abstracts the user-related operations
type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userID string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userID string, req *v1.UpdateProfileRequest) error
}

// NewUserService create a new UserService instance
func NewUserService(service *Service, userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

type userService struct {
	*Service
	userRepo repository.UserRepository
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check username
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return ErrDatabaseQuery
	}
	if err == nil && user != nil {
		return ErrUserYetExist
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Generate user ID
	userID, err := s.sid.GenString()
	if err != nil {
		return err
	}
	user = &model.User{
		UserID:   userID,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// TODO: Move transaction to repository layer
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		return s.userRepo.Create(ctx, user)
	})
	return err
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}
	token, err := s.jwt.GenToken(user.UserID, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:   user.UserID,
		Nickname: user.Nickname,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Nickname = req.Nickname
	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}
