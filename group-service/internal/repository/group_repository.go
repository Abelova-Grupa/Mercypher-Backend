package repository

import (
	"context"
	"errors"

	"github.com/Abelova-Grupa/Mercypher/group-service/internal/model"
	"gorm.io/gorm"
)

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *groupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) CreateGroup(
	ctx context.Context,
	group *model.Group,
	owner *model.GroupMember,
) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		groupModel := model.Group{
			ID:        group.ID,
			Name:      group.Name,
			OwnerID:   group.OwnerID,
			CreatedAt: group.CreatedAt,
		}

		if err := tx.Create(&groupModel).Error; err != nil {
			return err
		}

		memberModel := model.GroupMember{
			GroupID:  owner.GroupID,
			UserID:   owner.UserID,
			Role:     owner.Role,
			JoinedAt: owner.JoinedAt,
		}

		if err := tx.Create(&memberModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *groupRepository) GetMember(
	ctx context.Context,
	groupID, userID string,
) (*model.GroupMember, error) {

	var member model.GroupMember

	err := r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&member).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("Record ain't found.")
	}

	if err != nil {
		return nil, err
	}

	return &model.GroupMember{
		GroupID:  member.GroupID,
		UserID:   member.UserID,
		Role:     model.GroupRole(member.Role),
		JoinedAt: member.JoinedAt,
	}, nil
}

func (r *groupRepository) AddMember(
	ctx context.Context,
	member *model.GroupMember,
) error {

	memberModel := model.GroupMember{
		GroupID:  member.GroupID,
		UserID:   member.UserID,
		Role:     member.Role,
		JoinedAt: member.JoinedAt,
	}

	return r.db.WithContext(ctx).Create(&memberModel).Error
}

func (r *groupRepository) RemoveMember(
	ctx context.Context,
	groupID, userID string,
) error {

	return r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.GroupMember{}).Error
}

func (r *groupRepository) GetGroup(
	ctx context.Context,
	groupID string,
) (*model.Group, error) {

	var g model.Group

	err := r.db.WithContext(ctx).
		First(&g, "id = ?", groupID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &model.Group{
		ID:        g.ID,
		Name:      g.Name,
		OwnerID:   g.OwnerID,
		CreatedAt: g.CreatedAt,
	}, nil
}

func (r *groupRepository) DeleteGroup(
	ctx context.Context,
	groupID string,
) error {

	return r.db.WithContext(ctx).
		Delete(&model.Group{}, "id = ?", groupID).
		Error
}

func (r *groupRepository) GetGroupMembers(
	ctx context.Context,
	groupID string,
) ([]*model.GroupMember, error) {

	var members []model.GroupMember

	err := r.db.WithContext(ctx).
		Where("group_id = ?", groupID).
		Find(&members).Error

	if err != nil {
		return nil, err
	}

	result := make([]*model.GroupMember, 0, len(members))

	for _, m := range members {
		result = append(result, &model.GroupMember{
			GroupID:  m.GroupID,
			UserID:   m.UserID,
			Role:     model.GroupRole(m.Role),
			JoinedAt: m.JoinedAt,
		})
	}

	return result, nil
}

func (r *groupRepository) GetUserGroups(
	ctx context.Context,
	userID string,
) ([]*model.Group, error) {

	var groups []model.Group

	err := r.db.WithContext(ctx).
		Table("group_service.groups").
		Select("groups.*").
		Joins("JOIN group_service.group_members gm ON gm.group_id = groups.id").
		Where("gm.user_id = ?", userID).
		Find(&groups).Error

	if err != nil {
		return nil, err
	}

	result := make([]*model.Group, 0, len(groups))

	for _, g := range groups {
		result = append(result, &model.Group{
			ID:        g.ID,
			Name:      g.Name,
			OwnerID:   g.OwnerID,
			CreatedAt: g.CreatedAt,
		})
	}

	return result, nil
}

func (r *groupRepository) UpdateMemberRole(
	ctx context.Context,
	groupID, userID string,
	role model.GroupRole,
) error {

	return r.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("role", int16(role)).
		Error
}

func (r *groupRepository) UpdateGroup(
	ctx context.Context,
	group *model.Group,
) error {
	return r.db.WithContext(ctx).
		Model(&model.Group{}).
		Where("id = ?", group.ID).
		Update("name", group.Name).Error
}