package service

import (
	"context"
	"errors"
	"testing"

	pb "github.com/Abelova-Grupa/Mercypher/user-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define mock user repository
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	createError := args.Get(0)
	if createError != nil {
		return args.Error(0)
	}
	return nil
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	user := args.Get(0)
	err := args.Get(1)
	if user == nil || err != nil {
		return nil, args.Error(1)
	}
	// Casting
	return user.(*models.User), nil
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	user := args.Get(0)
	err := args.Get(1)
	if user == nil || err != nil {
		return nil, args.Error(1)
	}
	return user.(*models.User), nil
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	updateError := args.Get(0)
	if updateError != nil {
		return args.Error(0)
	}
	return nil
}

func (m *MockUserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	user := args.Get(0)
	err := args.Get(1)
	if user == nil || err != nil {
		return nil, args.Error(1)
	}
	return user.(*models.User), nil
}

func (m *MockUserRepo) Login(ctx context.Context, username string, password string) bool {
	args := m.Called(ctx, username, password)
	loggedIn := args.Get(0)
	return loggedIn.(bool)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)

	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").
		Return(&models.User{Username: "testuser"}, nil)

	_, err := userService.Register(context.Background(), &pb.User{
		Username: "testuser",
		Password: "password123",
	})
	assert.EqualError(t, err, "username already exists")
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)

	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").
		Return(nil, errors.New("not found"))

	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Username == "testuser" &&
			u.Email == "testuser@gmail.com" &&
			len(u.PasswordHash) > 0
	})).Return(nil)

	user, err := userService.Register(context.Background(), &pb.User{
		Username: "testuser",
		Email:    "testuser@gmail.com",
		Password: "password123",
	})
	assert.NotEqual(t, user, nil)
	assert.Equal(t, err, nil)
}
