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
	"github.com/Abelova-Grupa/Mercypher/user-service/token"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
)

var (
	ErrInvalidParams = errors.New("parameters are invalid")
	ErrInvalidEnvVars = errors.New("invalid env variables")
)

var (
	sessionDuration = 1440 * time.Minute
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, registerUserRequestPb *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {

	g, groupCtx := errgroup.WithContext(ctx)
	var hashed []byte

	g.Go(func() error {
		if _, err := s.repo.GetUserByUsername(groupCtx, registerUserRequestPb.GetUsername()); err == nil {
			return errors.New("username already exists")
		}
		return nil
	})

	g.Go(func() error {
		var err error
		hashed, err = bcrypt.GenerateFromPassword([]byte(registerUserRequestPb.GetPassword()), bcrypt.DefaultCost) 
		return err
	})

	// TODO: Rename password hash, not good variable name
	authCode := ""
	for i := 0; i < 5; i++ {
		authCode += fmt.Sprintf("%d",rand.Intn(10))
	}
	
	if err := g.Wait(); err != nil {
		return nil, err
	}

	user := &models.User{
		Username: registerUserRequestPb.GetUsername(),
		Email: registerUserRequestPb.GetEmail(),
		CreatedAt: registerUserRequestPb.GetCreatedAt().AsTime().Unix(),
		PasswordHash: string(hashed),
		Validated: false,
		AuthCode: authCode,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if err := s.SendEmail(user.Email,user.Username,authCode); err != nil {
		// User is created but has validated: false flag
		return nil, err
	}

	return &userpb.RegisterUserResponse{Username: user.Username,
		Email: user.Email,
		AuthCode: user.AuthCode},
		nil
}

func (s *UserService) ValidateAccount(ctx context.Context, validateRequest *userpb.ValidateUserAccountRequest) error {
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

func (u *UserService) CreateToken(ctx context.Context, username string, duration time.Duration) (string,error) {
	jwtMaker := token.JWTMaker{}
	token, _, err := jwtMaker.CreateToken(username, duration)
	if token == "" || err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserService) VerifyToken(ctx context.Context, verifyTokenRequestPb *pb.VerifyTokenRequest) (bool, error) {
	jwtMaker := token.JWTMaker{}
	payload, err := jwtMaker.VerifyToken(verifyTokenRequestPb.Token)
	if payload == nil || err != nil {
		return false, err
	}
	return true, nil
}

