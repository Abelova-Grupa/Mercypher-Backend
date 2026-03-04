package repository

import (
	"context"
	"github.com/Abelova-Grupa/Mercypher/group-service/internal/model"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *model.Group, owner *model.GroupMember) error
	GetGroup(ctx context.Context, groupID string) (*model.Group, error)
	DeleteGroup(ctx context.Context, groupID string) error
	UpdateGroup(ctx context.Context, group *model.Group) error

	AddMember(ctx context.Context, member *model.GroupMember) error
	RemoveMember(ctx context.Context, groupID, userID string) error
	GetMember(ctx context.Context, groupID, userID string) (*model.GroupMember, error)

	GetGroupMembers(ctx context.Context, groupID string) ([]*model.GroupMember, error)
	GetUserGroups(ctx context.Context, userID string) ([]*model.Group, error)
	UpdateMemberRole(ctx context.Context, groupID, userID string, role model.GroupRole) error
}