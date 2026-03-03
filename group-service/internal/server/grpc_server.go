package server

import (
	"context"
	"errors"
	"github.com/Abelova-Grupa/Mercypher/group-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/group-service/internal/repository"
	grouppb "github.com/Abelova-Grupa/Mercypher/proto/group"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GroupServer struct {
	grouppb.UnimplementedGroupServiceServer
	repo repository.GroupRepository
}

func NewGroupServer(repo repository.GroupRepository) *GroupServer {
	return &GroupServer{repo: repo}
}

func (s *GroupServer) requireOwner(ctx context.Context, groupID, userID string) error {
	member, err := s.repo.GetMember(ctx, groupID, userID)
	if err != nil {
		return status.Error(codes.PermissionDenied, "not a group member")
	}
	if !member.Role.IsOwner() {
		return status.Error(codes.PermissionDenied, "only owner allowed")
	}
	return nil
}

func toProtoRole(r model.GroupRole) grouppb.GroupRole {
	switch r {
	case model.RoleOwner:
		return grouppb.GroupRole_GROUP_ROLE_OWNER
	case model.RoleMember:
		return grouppb.GroupRole_GROUP_ROLE_MEMBER
	default:
		return grouppb.GroupRole_GROUP_ROLE_UNSPECIFIED
	}
}

func toModelRole(r grouppb.GroupRole) model.GroupRole {
	switch r {
	case grouppb.GroupRole_GROUP_ROLE_OWNER:
		return model.RoleOwner
	case grouppb.GroupRole_GROUP_ROLE_MEMBER:
		return model.RoleMember
	default:
		return model.RoleUnspecified
	}
}

func (s *GroupServer) CreateGroup(
	ctx context.Context,
	req *grouppb.CreateGroupRequest,
) (*grouppb.CreateGroupResponse, error) {

	if req.GetName() == "" || req.GetCreatorId() == "" {
		return nil, status.Error(codes.InvalidArgument, "name and creator_id required")
	}

	groupID := uuid.NewString()

	group := model.NewGroup(
		groupID,
		req.GetName(),
		req.GetCreatorId(),
	)

	owner := model.NewOwnerMember(
		groupID,
		req.GetCreatorId(),
	)

	if err := s.repo.CreateGroup(ctx, group, owner); err != nil {
		return nil, status.Errorf(codes.Internal, "create group failed: %v", err)
	}

	return &grouppb.CreateGroupResponse{
		Group: &grouppb.Group{
			Id:        group.ID,
			Name:      group.Name,
			OwnerId:   group.OwnerID,
			CreatedAt: timestamppb.New(group.CreatedAt),
		},
	}, nil
}

func (s *GroupServer) DeleteGroup(
	ctx context.Context,
	req *grouppb.DeleteGroupRequest,
) (*emptypb.Empty, error) {

	if err := s.requireOwner(ctx, req.GetGroupId(), req.GetRequesterId()); err != nil {
		return nil, err
	}

	if err := s.repo.DeleteGroup(ctx, req.GetGroupId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *GroupServer) UpdateGroup(
	ctx context.Context,
	req *grouppb.UpdateGroupRequest,
) (*grouppb.Group, error) {

	if err := s.requireOwner(ctx, req.GetGroupId(), req.GetRequesterId()); err != nil {
		return nil, err
	}

	group, err := s.repo.GetGroup(ctx, req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "group not found")
	}

	group.Name = req.GetNewName()

	if err := s.repo.UpdateGroup(ctx, group); err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}

	return &grouppb.Group{
		Id:        group.ID,
		Name:      group.Name,
		OwnerId:   group.OwnerID,
		CreatedAt: timestamppb.New(group.CreatedAt),
	}, nil
}

func (s *GroupServer) AddMember(
	ctx context.Context,
	req *grouppb.AddMemberRequest,
) (*emptypb.Empty, error) {

	if err := s.requireOwner(ctx, req.GetGroupId(), req.GetRequesterId()); err != nil {
		return nil, err
	}

	member := &model.GroupMember{
		GroupID:  req.GetGroupId(),
		UserID:   req.GetUserId(),
		Role:     model.RoleMember,
		JoinedAt: timestamppb.Now().AsTime(),
	}

	if err := s.repo.AddMember(ctx, member); err != nil {
		return nil, status.Errorf(codes.Internal, "add member failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *GroupServer) RemoveMember(
	ctx context.Context,
	req *grouppb.RemoveMemberRequest,
) (*emptypb.Empty, error) {

	if err := s.requireOwner(ctx, req.GetGroupId(), req.GetRequesterId()); err != nil {
		return nil, err
	}

	if err := s.repo.RemoveMember(ctx, req.GetGroupId(), req.GetUserId()); err != nil {
		return nil, status.Errorf(codes.Internal, "remove failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *GroupServer) ChangeMemberRole(
	ctx context.Context,
	req *grouppb.ChangeMemberRoleRequest,
) (*emptypb.Empty, error) {

	if err := s.requireOwner(ctx, req.GetGroupId(), req.GetRequesterId()); err != nil {
		return nil, err
	}

	role := toModelRole(req.GetNewRole())

	if role == model.RoleUnspecified {
		return nil, status.Error(codes.InvalidArgument, "invalid role")
	}

	if err := s.repo.UpdateMemberRole(ctx, req.GetGroupId(), req.GetUserId(), role); err != nil {
		return nil, status.Errorf(codes.Internal, "role change failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *GroupServer) GetGroup(
	ctx context.Context,
	req *grouppb.GetGroupRequest,
) (*grouppb.Group, error) {

	group, err := s.repo.GetGroup(ctx, req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "group not found")
	}

	return &grouppb.Group{
		Id:        group.ID,
		Name:      group.Name,
		OwnerId:   group.OwnerID,
		CreatedAt: timestamppb.New(group.CreatedAt),
	}, nil
}

func (s *GroupServer) GetGroupMembers(
	ctx context.Context,
	req *grouppb.GetGroupMembersRequest,
) (*grouppb.GetGroupMembersResponse, error) {

	members, err := s.repo.GetGroupMembers(ctx, req.GetGroupId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fetch failed: %v", err)
	}

	var pbMembers []*grouppb.GroupMember

	for _, m := range members {
		pbMembers = append(pbMembers, &grouppb.GroupMember{
			GroupId:  m.GroupID,
			UserId:   m.UserID,
			Role:     toProtoRole(m.Role),
			JoinedAt: timestamppb.New(m.JoinedAt),
		})
	}

	return &grouppb.GetGroupMembersResponse{
		Members: pbMembers,
	}, nil
}

func (s *GroupServer) GetUserGroups(
	ctx context.Context,
	req *grouppb.GetUserGroupsRequest,
) (*grouppb.GetUserGroupsResponse, error) {

	groups, err := s.repo.GetUserGroups(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fetch failed: %v", err)
	}

	var pbGroups []*grouppb.Group

	for _, g := range groups {
		pbGroups = append(pbGroups, &grouppb.Group{
			Id:        g.ID,
			Name:      g.Name,
			OwnerId:   g.OwnerID,
			CreatedAt: timestamppb.New(g.CreatedAt),
		})
	}

	return &grouppb.GetUserGroupsResponse{
		Groups: pbGroups,
	}, nil
}

func (s *GroupServer) IsMember(
	ctx context.Context,
	req *grouppb.IsMemberRequest,
) (*grouppb.IsMemberResponse, error) {

	member, err := s.repo.GetMember(ctx, req.GetGroupId(), req.GetUserId())
	if err != nil {
		if errors.Is(err, status.Error(codes.PermissionDenied, "")) {
			return &grouppb.IsMemberResponse{
				IsMember: false,
			}, nil
		}
		return &grouppb.IsMemberResponse{
			IsMember: false,
		}, nil
	}

	return &grouppb.IsMemberResponse{
		IsMember: true,
		Role:     toProtoRole(member.Role),
	}, nil
}