package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/email"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/token"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/worker"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrInvalidParams  = errors.New("parameters are invalid")
	ErrInvalidEnvVars = errors.New("invalid env variables")
)

type UserService struct {
	repo            repository.UserRepository
	taskDistributor worker.TaskDistributor
	db              *gorm.DB
}

type RegisterUserInput struct {
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}

type RegisterUserResponse struct {
	Username  string
	Email     string
	CreatedAt time.Time
}

type LoginUserInput struct {
	Username string
	Password string
}

type TokenInput struct {
	Token string
}

type ValidateAccountInput struct {
	Username string
	AuthCode string
}

type SendEmailInput struct {
	Username string
	Email    string
	AuthCode string
}

type CreateTokenInput struct {
	Username string
	Duration time.Duration
}

type DecodeAccessTokenInput struct {
	Token string
}

type CreateContactInput struct {
	Username    string
	ContactName string
}

type DeleteContactInput struct {
	Username    string
	ContactName string
}

type GetContactsInput struct {
	Username       string
	SearchCriteria string
}

func NewUserService(db *gorm.DB, repo repository.UserRepository) *UserService {
	redisOpt := asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASS"),
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	return &UserService{repo: repo, taskDistributor: taskDistributor, db: db}
}

func (s *UserService) Register(ctx context.Context, input RegisterUserInput) (*RegisterUserResponse, error) {
	g, groupCtx := errgroup.WithContext(ctx)
	var hashed []byte

	g.Go(func() error {
		if user, _ := s.repo.GetUserByUsername(groupCtx, input.Username); user != nil {
			return errors.New("username already exists")
		}
		return nil
	})

	g.Go(func() error {
		var err error
		hashed, err = bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		return err
	})

	authCode := ""
	for i := 0; i < 5; i++ {
		authCode += fmt.Sprintf("%d", rand.Intn(10))
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		CreatedAt:    input.CreatedAt,
		PasswordHash: string(hashed),
		Validated:    false,
		AuthCode:     authCode,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		repoWithTx := s.repo.WithTx(tx)

		if err := repoWithTx.CreateUser(ctx, user); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	payload := &email.EmailPayload{
		Username: user.Username,
		ToEmail:  user.Email,
		AuthCode: user.AuthCode,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.ProcessIn(5 * time.Second),
	}

	if err := s.taskDistributor.DistributeTaskSendVerifyEmail(ctx, payload, opts...); err != nil {
		return nil, err
	}

	return &RegisterUserResponse{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil

}

func (s *UserService) Login(ctx context.Context, input LoginUserInput) (bool, error) {
	isLoggedIn := s.repo.Login(ctx, input.Username, input.Password)
	return isLoggedIn, nil
}

func (s *UserService) ValidateAccount(ctx context.Context, input ValidateAccountInput) error {
	return s.repo.ValidateAccount(ctx, input.Username, input.AuthCode)
}

// TODO: Think about adding context here for timeout reasons
func (u *UserService) CreateToken(ctx context.Context, input CreateTokenInput) (string, error) {
	jwtMaker := token.JWTMaker{}
	token, _, err := jwtMaker.CreateToken(input.Username, input.Duration)
	if token == "" || err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserService) VerifyToken(ctx context.Context, tokenRequest TokenInput) (bool, error) {
	jwtMaker := token.JWTMaker{}
	payload, err := jwtMaker.VerifyToken(tokenRequest.Token)
	if payload == nil || err != nil {
		return false, err
	}
	return true, nil
}

func (s *UserService) DecodeAccessToken(ctx context.Context, input DecodeAccessTokenInput) (string, error) {
	jwtMaker := token.JWTMaker{}
	payload, err := jwtMaker.VerifyToken(input.Token)
	if payload == nil || err != nil {
		return "", err
	}
	return payload.UserID, nil
}

func (s *UserService) CreateContact(ctx context.Context, input *CreateContactInput) (*models.Contact, error) {
	contact, err := s.repo.CreateContact(ctx, input.Username, input.ContactName)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (s *UserService) DeleteContact(ctx context.Context, input *DeleteContactInput) error {

	res := s.db.Delete(&models.Contact{}, " username = ? AND contact_name = ?", input.Username, input.ContactName)
	if res.Error != nil {
		return status.Error(codes.Internal, "database error")
	}
	if res.RowsAffected == 0 {
		return status.Error(codes.NotFound, "contact doesn't exist")
	}
	return nil
}

func (s *UserService) GetContactsStream(ctx context.Context, input *GetContactsInput, handleSend func(models.Contact) error) error {
	chanContact, chanErr := s.repo.GetContactsCursor(ctx, input.Username, input.SearchCriteria)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-chanErr:
			if err != nil {
				return err
			}
			return nil
		case c, ok := <-chanContact:
			if !ok {
				return fmt.Errorf("unable to retreive contact %v", c)
			}
			if err := handleSend(c); err != nil {
				return err
			}
		}
	}

}
