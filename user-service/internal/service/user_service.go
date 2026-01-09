package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	pb "github.com/Abelova-Grupa/Mercypher/proto/user"
	userpb "github.com/Abelova-Grupa/Mercypher/proto/user"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrInvalidParams = errors.New("parameters are invalid")
	ErrInvalidEnvVars = errors.New("invalid env variables")
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

	authCode := ""
	for i := 0; i < 5; i++ {
		authCode += fmt.Sprintf("%d",rand.Intn(10))
	}
	user.AuthCode = authCode

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if err := s.SendEmail(user.Email,user.Username,authCode); err != nil {
		//TODO: Implement account deletion
		return nil, err
	}

	return convertUserToPb(user), nil
}

func (s *UserService) ValidateAccount(ctx context.Context, validateRequest *userpb.ValidateAccountRequest) error {
	if validateRequest == nil || validateRequest.Username == "" || validateRequest.AuthCode == "" {
		return ErrInvalidParams
	}
	return s.repo.ValidateAccount(ctx, validateRequest.Username, validateRequest.AuthCode)
}

func (s *UserService) SendEmail(toEmail string, username string, authCode string) error {
	appPass := os.Getenv("EMAIL_APP_PASS")
	fromEmail := os.Getenv("EMAIL_CLIENT")
	smtpServer := os.Getenv("SMTP_HOST")
	tlsPort := os.Getenv("SSL_SMTP_PORT")

	if appPass == "" || fromEmail == "" || smtpServer == "" || tlsPort == "" {
		return ErrInvalidEnvVars
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpServer,
	}

	auth := smtp.PlainAuth("", fromEmail, appPass, smtpServer)
	conn, err := tls.Dial("tcp", smtpServer+":"+tlsPort, tlsConfig)
	if err != nil {
		return err
	}

	headers := make(map[string] string)
	headers["From"] = fromEmail
	headers["To"] = toEmail
	headers["Subject"] = "Verify your mercypher account"

	messageBody := "Your verification code is " + authCode + ". You have 15 minutes to active your account, otherwise it will be deleted"
	message := ""
	for k,v := range headers {
		message += fmt.Sprintf("%s: %s\r\n",k,v)
	}
	message += messageBody
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(fromEmail); err != nil {
		return err
	}

	if err := client.Rcpt(toEmail); err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}

	_, err = wc.Write([]byte(message))
	if err != nil {
		return nil
	}
	wc.Close()
	return nil
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
