package server

import (
	"context"
	"errors"
	"time"

	pb "github.com/Abelova-Grupa/Mercypher/session-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/services"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type grpcServer struct {
	sessionDB      *gorm.DB
	sessionRepo    repository.SessionRepository
	sessionService services.SessionService
	pb.UnsafeSessionServiceServer
}

func NewGrpcServer(db *gorm.DB) *grpcServer {
	repo := repository.NewSessionRepository(db)
	service := services.NewSessionService(repo)
	return &grpcServer{
		sessionDB:      db,
		sessionRepo:    repo,
		sessionService: *service,
	}
}

func (s *grpcServer) GetUserLocation(ctx context.Context, userID *pb.UserID) (*pb.UserLocation, error) {

	userLocation, err := s.sessionRepo.GetUserLocationByUserID(ctx, userID.UserID)
	if err != nil {
		return &pb.UserLocation{}, err
	}

	return &pb.UserLocation{
		UserID:    userLocation.UserID,
		APIAdress: userLocation.ApiIP,
	}, nil
}

func (s *grpcServer) UpdateUserLocation(ctx context.Context, userLoc *pb.UserLocation) (*pb.UserLocation, error) {
	userLocation := models.UserLocation{
		UserID: userLoc.UserID,
		ApiIP:  userLoc.APIAdress,
	}
	// If the userID doesnt exist it will create a new UserLocation, otherwise it will update existing UserLocation
	err := s.sessionRepo.UpdateUserLocation(ctx, &userLocation)
	if err != nil {
		return &pb.UserLocation{}, errors.New("unable to update user location")
	}

	updatedUserLocation, err := s.sessionRepo.GetUserLocationByUserID(ctx, userLoc.UserID)
	if err != nil {
		return &pb.UserLocation{}, errors.New("unable to retreive updated user location")
	}
	return &pb.UserLocation{
		UserID:    updatedUserLocation.UserID,
		APIAdress: updatedUserLocation.ApiIP,
	}, nil
}

func (s *grpcServer) GetLastSeen(ctx context.Context, userID *pb.UserID) (*pb.LastSeen, error) {
	lastSeen, err := s.sessionRepo.GetLastSeenByUserID(ctx, userID.UserID)
	if err != nil {
		return &pb.LastSeen{}, errors.New("unable to retreive last seen info")
	}
	return &pb.LastSeen{
		UserID:   lastSeen.UserID,
		LastSeen: timestamppb.New(time.Unix(lastSeen.LastSeen, 0)),
	}, nil
}

func (s *grpcServer) UpdateLastSeen(ctx context.Context, lastSeen *pb.LastSeen) (*pb.LastSeen, error) {
	ls := models.LastSeenSession{
		UserID:   lastSeen.UserID,
		LastSeen: lastSeen.LastSeen.AsTime().Unix(),
	}
	err := s.sessionRepo.UpdateLastSeen(ctx, &ls)
	if err != nil {
		return &pb.LastSeen{}, errors.New("unable to update last seen info")
	}

	lastSeenUpdated, err := s.sessionRepo.GetLastSeenByUserID(ctx, lastSeen.UserID)
	if err != nil {
		return &pb.LastSeen{}, errors.New("unable to retreive last seen info")
	}
	return &pb.LastSeen{
		UserID:   lastSeenUpdated.UserID,
		LastSeen: timestamppb.New(time.Unix(lastSeenUpdated.LastSeen, 0)),
	}, nil

}
