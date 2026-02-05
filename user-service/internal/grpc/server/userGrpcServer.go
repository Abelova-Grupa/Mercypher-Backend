package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	sessionClient "github.com/Abelova-Grupa/Mercypher/session-service/external/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	sessionpb "github.com/Abelova-Grupa/Mercypher/proto/session"
	userpb "github.com/Abelova-Grupa/Mercypher/proto/user"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/service"
	"gorm.io/gorm"
)

type GrpcServer struct {
	userDB      *gorm.DB
	userRepo    repository.UserRepository
	userService service.UserService
	userpb.UnsafeUserServiceServer
	sessionClient sessionClient.GrpcClient
}

func NewGrpcServer(db *gorm.DB) *GrpcServer {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(db, repo)

	var sessionUrl string
	sessionPort := os.Getenv("SESSION_SERVICE_PORT")
	if os.Getenv("ENVIRONMENT") == "" {
		sessionUrl = fmt.Sprintf("localhost:%s", sessionPort)
	} else {
		sessionUrl = fmt.Sprintf("session-service:%s", sessionPort)
	}

	grpcClient, _ := sessionClient.NewGrpcClient(sessionUrl)
	return &GrpcServer{
		userDB:        db,
		userRepo:      repo,
		userService:   *service,
		sessionClient: *grpcClient,
	}
}

// Should only create a user not a session
func (g *GrpcServer) RegisterUser(ctx context.Context, registerRequestPb *userpb.RegisterUserRequest) (*userpb.RegisterUserResponse, error) {
	if registerRequestPb == nil || registerRequestPb.Username == "" || registerRequestPb.Email == "" || registerRequestPb.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email and password are required for registration")
	}
	res, err := g.userService.Register(ctx, service.RegisterUserInput{
		Username: registerRequestPb.Username,
		Email:    registerRequestPb.Email,
		Password: registerRequestPb.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	return &userpb.RegisterUserResponse{
		Username: res.Username,
		Email:    res.Email,
	}, nil
}

func (g *GrpcServer) LoginUser(ctx context.Context, loginRequest *userpb.LoginUserRequest) (*userpb.LoginUserResponse, error) {
	if loginRequest == nil || loginRequest.Username == "" || loginRequest.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required for login")
	}

	username := loginRequest.Username
	passedToken := loginRequest.Token
	password := loginRequest.Password

	if passedToken != "" {
		verified, _ := g.userService.VerifyToken(ctx, service.TokenInput{
			Token: passedToken,
		})
		if verified {
			return &userpb.LoginUserResponse{Username: username, AccessToken: passedToken}, nil
		} else {
			log.Print("Token is invalid, continue with credential checking")
		}
	}

	log.Println("Checking user credentials")
	isLoggedIn, _ := g.userService.Login(ctx, service.LoginUserInput{Username: username, Password: password})
	if !isLoggedIn {
		return nil, errors.New("Authentification failed")
	}
	log.Println("Successful authentication creating session...")
	var token string
	var err error
	if token, err = g.userService.CreateToken(ctx, service.CreateTokenInput{Username: username, Duration: 24 * time.Hour}); err != nil {
		return nil, fmt.Errorf("Failed to create auth token for user %v : %v", username, err)
	}
	_, err = g.sessionClient.Connect(ctx, &sessionpb.ConnectRequest{Username: username})
	if err != nil {
		return nil, fmt.Errorf("Failed session sign in for user %v : %v ", username, err)
	}

	log.Print("Succesfull login")
	return &userpb.LoginUserResponse{Username: username, AccessToken: token}, nil

}

func (g *GrpcServer) LogoutUser(ctx context.Context, logoutRequest *userpb.LogoutUserRequest) (*emptypb.Empty, error) {
	if logoutRequest.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username cannot be nil for logout")
	}
	usernamePb := &sessionpb.DisconnectRequest{Username: logoutRequest.Username}
	if _, err := g.sessionClient.Disconnect(ctx, usernamePb); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GrpcServer) ValidateUserAccount(ctx context.Context, validateRequest *userpb.ValidateUserAccountRequest) (*emptypb.Empty, error) {
	if validateRequest == nil || validateRequest.Username == "" || validateRequest.AuthCode == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments for account validation")
	}
	if err := g.userService.ValidateAccount(ctx,
		service.ValidateAccountInput{Username: validateRequest.Username, AuthCode: validateRequest.AuthCode}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GrpcServer) VerifyToken(ctx context.Context, verifyTokenRequest *userpb.VerifyTokenRequest) (*wrapperspb.BoolValue, error) {
	if verifyTokenRequest == nil || verifyTokenRequest.Token == "" {
		return wrapperspb.Bool(false), status.Error(codes.InvalidArgument, "cannot verify empty token")
	}

	if valid, err := g.userService.VerifyToken(ctx, service.TokenInput{Token: verifyTokenRequest.Token}); !valid || err != nil {
		return wrapperspb.Bool(false), err
	}
	return wrapperspb.Bool(true), nil
}

func (g *GrpcServer) DecodeAccessToken(ctx context.Context, decodeRequest *userpb.DecodeAccessTokenRequest) (*userpb.DecodeAccessTokenResponse, error) {
	if decodeRequest == nil || decodeRequest.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid parameters")
	}
	username, err := g.userService.DecodeAccessToken(ctx, service.DecodeAccessTokenInput{Token: decodeRequest.Token})
	if err != nil {
		return nil, status.Error(codes.NotFound, "couldn't return access token payload")
	}
	return &userpb.DecodeAccessTokenResponse{Username: username}, nil
}

func (g *GrpcServer) CreateContact(ctx context.Context, contactRequest *userpb.CreateContactRequest) (*userpb.CreateContactResponse, error) {
	if contactRequest == nil || contactRequest.Username == "" || contactRequest.ContactName == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments for contact creation")
	}
	contactInput := &service.CreateContactInput{
		Username:    contactRequest.Username,
		ContactName: contactRequest.ContactName,
		Nickname: contactRequest.Nickname,
	}

	contact, err := g.userService.CreateContact(ctx, contactInput)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "system couldn't create a contact")
	}
	contactRes := &userpb.CreateContactResponse{
		Username:    contact.Username,
		ContactName: contact.ContactName,
		Nickname: contact.Nickname,
		CreatedAt:   timestamppb.New(contact.CreatedAt),
		
	}

	return contactRes, nil
}

func (g *GrpcServer) DeleteContact(ctx context.Context, contactRequest *userpb.DeleteContactRequest) (*emptypb.Empty, error) {
	if contactRequest == nil || contactRequest.Username == "" || contactRequest.ContactName == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments for contact deletion")
	}

	contactInput := &service.DeleteContactInput{
		Username:    contactRequest.Username,
		ContactName: contactRequest.ContactName,
	}

	err := g.userService.DeleteContact(ctx, contactInput)
	return &emptypb.Empty{}, err
}

func (g *GrpcServer) GetContacts(contactRequest *userpb.GetContactsRequest, stream userpb.UserService_GetContactsServer) error {
	if contactRequest == nil || contactRequest.Username == "" {
		return status.Error(codes.InvalidArgument, "invalid arguments for contact retreival")
	}

	ctx := stream.Context()
	return g.userService.GetContactsStream(ctx, &service.GetContactsInput{Username: contactRequest.Username, SearchCriteria: contactRequest.SearchCriteria}, func(c models.Contact) error {
		return stream.Send(&userpb.GetContactResponse{
			Username:    c.Username,
			ContactName: c.ContactName,
		})
	})

}

func (g *GrpcServer) UpdateContact(ctx context.Context, contactRequest *userpb.UpdateContactRequest) (*userpb.UpdateContactResponse, error) {
	if contactRequest == nil || contactRequest.Username == "" || contactRequest.ContactName == "" {
		return nil, status.Error(codes.InvalidArgument,"invalid arguments for contact update")
	}

	contactInput := &service.UpdateContactInput{
		Username:    contactRequest.Username,
		ContactName: contactRequest.ContactName,
		Nickname: contactRequest.Nickname,
	}

	contact, err := g.userService.UpdateContact(ctx, contactInput)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "system couldn't update a contact")
	}
	contactRes := &userpb.UpdateContactResponse{
		Username:    contact.Username,
		ContactName: contact.ContactName,
		Nickname:   contact.Nickname,
	}
	return contactRes, nil
}
