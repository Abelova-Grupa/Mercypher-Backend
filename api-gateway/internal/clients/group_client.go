package clients

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grouppb "github.com/Abelova-Grupa/Mercypher/proto/group"
)

type GroupClient struct {
	conn   *grpc.ClientConn
	client grouppb.GroupServiceClient
}

func NewGroupClient(address string) (*GroupClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, errors.New("connection refused: nil")
	}

	client := grouppb.NewGroupServiceClient(conn)

	return &GroupClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *GroupClient) Close() error {
	return c.conn.Close()
}

func (c *GroupClient) CreateGroup(ctx context.Context, name, creatorID string) (*grouppb.Group, error) {
	req := &grouppb.CreateGroupRequest{
		Name:      name,
		CreatorId: creatorID,
	}

	resp, err := c.client.CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Group, nil
}

func (c *GroupClient) GetGroup(ctx context.Context, groupID string) (*grouppb.Group, error) {
	req := &grouppb.GetGroupRequest{
		GroupId: groupID,
	}

	return c.client.GetGroup(ctx, req)
}

func (c *GroupClient) AddMember(ctx context.Context, groupID, requesterID, userID string) error {
	req := &grouppb.AddMemberRequest{
		GroupId:     groupID,
		RequesterId: requesterID,
		UserId:      userID,
	}

	_, err := c.client.AddMember(ctx, req)
	return err
}

func (c *GroupClient) IsMember(ctx context.Context, groupID, userID string) (bool, grouppb.GroupRole, error) {
	req := &grouppb.IsMemberRequest{
		GroupId: groupID,
		UserId:  userID,
	}

	resp, err := c.client.IsMember(ctx, req)
	if err != nil {
		return false, grouppb.GroupRole_GROUP_ROLE_UNSPECIFIED, err
	}

	return resp.IsMember, resp.Role, nil
}

func (c *GroupClient) GetUserGroups(ctx context.Context, userID string) ([]*grouppb.Group, error) {
	req := &grouppb.GetUserGroupsRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUserGroups(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Groups, nil
}

func (c *GroupClient) DeleteGroup(ctx context.Context, groupID, requesterID string) error {
	req := &grouppb.DeleteGroupRequest{
		GroupId:     groupID,
		RequesterId: requesterID,
	}

	_, err := c.client.DeleteGroup(ctx, req)
	return err
}

func (c *GroupClient) UpdateGroup(ctx context.Context, groupID, requesterID, newName string) (*grouppb.Group, error) {
	req := &grouppb.UpdateGroupRequest{
		GroupId:     groupID,
		RequesterId: requesterID,
		NewName:     newName,
	}

	return c.client.UpdateGroup(ctx, req)
}

func (c *GroupClient) RemoveMember(ctx context.Context, groupID, requesterID, userID string) error {
	req := &grouppb.RemoveMemberRequest{
		GroupId:     groupID,
		RequesterId: requesterID,
		UserId:      userID,
	}

	_, err := c.client.RemoveMember(ctx, req)
	return err
}

func (c *GroupClient) ChangeMemberRole(ctx context.Context, groupID, requesterID, userID string, newRole int32) error {
	req := &grouppb.ChangeMemberRoleRequest{
		GroupId:     groupID,
		RequesterId: requesterID,
		UserId:      userID,
		NewRole:     grouppb.GroupRole(newRole),
	}

	_, err := c.client.ChangeMemberRole(ctx, req)
	return err
}

func (c *GroupClient) GetGroupMembers(ctx context.Context, groupID string) ([]*grouppb.GroupMember, error) {
	req := &grouppb.GetGroupMembersRequest{
		GroupId: groupID,
	} 

	resp, err := c.client.GetGroupMembers(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Members, nil
}