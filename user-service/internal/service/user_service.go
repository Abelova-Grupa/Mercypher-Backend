package service

import (
	"context"
	"errors"
	"time"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidParams = errors.New("parameters are invalid")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	_, err := s.repo.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashed),
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (bool, error) {
	if username == "" || password == "" {
		return false, ErrInvalidParams
	}
	isLoggedIn := s.repo.Login(ctx, username, password)
	return isLoggedIn, nil
}
