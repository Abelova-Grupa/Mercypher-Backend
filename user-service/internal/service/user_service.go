package service

import (
	"context"
	"errors"

	pb "github.com/Abelova-Grupa/Mercypher/user-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *UserService) Register(ctx context.Context, userPb *pb.User) (*pb.User, error) {
	user := convertPbToUser(userPb)

	_, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	// TODO: Rename password hash, not good variable name
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Other fields are already stored in user struct
	user.ID = uuid.New().String()
	user.PasswordHash = string(hashed)

	// user := &models.User{
	// 	ID:           uuid.New().String(),
	// 	Username:     username,
	// 	Email:        email,
	// 	PasswordHash: string(hashed),
	// 	CreatedAt:    time.Now(),
	// }

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return convertUserToPb(user), nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (bool, error) {
	if username == "" || password == "" {
		return false, ErrInvalidParams
	}
	isLoggedIn := s.repo.Login(ctx, username, password)
	return isLoggedIn, nil
}

func convertPbToUser(userPb *pb.User) *models.User {
	return &models.User{
		ID:           userPb.GetID(),
		Username:     userPb.Username,
		Email:        userPb.Email,
		PasswordHash: userPb.GetPassword(),
		CreatedAt:    userPb.GetCreatedAt().AsTime(),
	}
}

func convertUserToPb(user *models.User) *pb.User {
	return &pb.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.PasswordHash,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
