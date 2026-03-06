package server

import (
	"context"
	"testing"
	"time"

	grouppb "github.com/Abelova-Grupa/Mercypher/proto/group"
	"github.com/Abelova-Grupa/Mercypher/group-service/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	createGroupFn func(ctx context.Context, g *model.Group, m *model.GroupMember) error
	getMemberFn   func(ctx context.Context, groupID, userID string) (*model.GroupMember, error)
}

func (m *mockRepo) CreateGroup(ctx context.Context, g *model.Group, mem *model.GroupMember) error {
	return m.createGroupFn(ctx, g, mem)
}

func (m *mockRepo) GetMember(ctx context.Context, groupID, userID string) (*model.GroupMember, error) {
	return m.getMemberFn(ctx, groupID, userID)
}

func (m *mockRepo) DeleteGroup(context.Context, string) error                         { return nil }
func (m *mockRepo) UpdateGroup(context.Context, *model.Group) error                   { return nil }
func (m *mockRepo) AddMember(context.Context, *model.GroupMember) error               { return nil }
func (m *mockRepo) RemoveMember(context.Context, string, string) error                { return nil }
func (m *mockRepo) GetGroup(context.Context, string) (*model.Group, error)            { return nil, nil }
func (m *mockRepo) GetGroupMembers(context.Context, string) ([]*model.GroupMember, error) {
	return nil, nil
}
func (m *mockRepo) GetUserGroups(context.Context, string) ([]*model.Group, error) {
	return nil, nil
}
func (m *mockRepo) UpdateMemberRole(context.Context, string, string, model.GroupRole) error {
	return nil
}

func TestCreateGroup_Success(t *testing.T) {
	mock := &mockRepo{
		createGroupFn: func(ctx context.Context, g *model.Group, m *model.GroupMember) error {
			assert.Equal(t, g.OwnerID, m.UserID)
			assert.Equal(t, model.RoleOwner, m.Role)
			return nil
		},
	}

	server := NewGroupServer(mock)

	resp, err := server.CreateGroup(context.Background(), &grouppb.CreateGroupRequest{
		Name:      "Test",
		CreatorId: "slavko",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Test", resp.Group.Name)
	assert.Equal(t, "slavko", resp.Group.OwnerId)
}

func TestDeleteGroup_PermissionDenied(t *testing.T) {
	mock := &mockRepo{
		getMemberFn: func(ctx context.Context, groupID, userID string) (*model.GroupMember, error) {
			return &model.GroupMember{
				GroupID:  groupID,
				UserID:   userID,
				Role:     model.RoleMember,
				JoinedAt: time.Now(),
			}, nil
		},
	}

	server := NewGroupServer(mock)

	_, err := server.DeleteGroup(context.Background(), &grouppb.DeleteGroupRequest{
		GroupId:    "group1",
		RequesterId: "slavko",
	})

	assert.Error(t, err)
}