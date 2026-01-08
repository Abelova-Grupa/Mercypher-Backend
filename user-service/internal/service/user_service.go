package service

import (
	"context"
	"crypto/tls"
	"errors"
	"net/smtp"
	"os"
	"time"

	pb "github.com/Abelova-Grupa/Mercypher/proto/user"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
	"github.com/gin-gonic/gin"
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
	user.PasswordHash = string(hashed)
	user.Validated = false

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return convertUserToPb(user), nil
}

func (s *UserService) ValidateAccount(ctx context.Context, username string, authCode string) error {
	if username == "" || authCode == "" {
		return ErrInvalidParams
	}
	return s.repo.ValidateAccount(ctx, username, authCode)
}

func (s *UserService) SendEmail(ctx *gin.Context, toEmail string, username string) error {
	appPass := os.Getenv("EMAIL_APP_PASS")
	fromEmail := os.Getenv("EMAIL_CLIENT")
	smtpServer := os.Getenv("SMTP_HOST")
	tlsPort := os.Getenv("TLS_SMTP_PORT")

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpServer,
	}

	auth := smtp.PlainAuth("", fromEmail, appPass, smtpServer)

	conn, err := tls.Dial("tcp", smtpServer+":"+tlsPort, tlsConfig)
	if err != nil {

	}

	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {

	}

	if err := client.Auth(auth); err != nil {

	}

	if err := client.Mail(fromEmail); err != nil {

	}

	if err := client.Rcpt(toEmail); err != nil {

	}

	wc, err := client.Data()
	if err != nil {

	}

	_, err = wc.Write([]byte(message))
	if err != nil {

	}

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
		Username:     userPb.Username,
		Email:        userPb.Email,
		PasswordHash: userPb.GetPassword(),
		CreatedAt:    userPb.GetCreatedAt().AsTime().Unix(),
	}
}

func convertUserToPb(user *models.User) *pb.User {
	return &pb.User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.PasswordHash,
		CreatedAt: timestamppb.New(time.Unix(user.CreatedAt, 0)),
	}
}
