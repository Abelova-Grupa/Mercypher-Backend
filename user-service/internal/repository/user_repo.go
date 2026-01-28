package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	Login(ctx context.Context, username string, password string) bool
	ValidateAccount(ctx context.Context, username string, authCode string) error
	CreateContact(ctx context.Context, user1 *models.User, user2 *models.User) (*models.Contact, error)
}

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepo{DB: db}
}

func (r *UserRepo) WithTx(tx *gorm.DB) UserRepository {
	return &UserRepo{DB: tx}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	result := r.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &user, result.Error
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	result := r.DB.WithContext(ctx).First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) Login(ctx context.Context, username string, password string) bool {
	var user models.User

	err := r.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if user.Validated == false {
		return false
	} else if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

func (r *UserRepo) ValidateAccount(ctx context.Context, username string, authCode string) error {
	var user models.User
	result := r.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	if user.AuthCode != authCode {
		return fmt.Errorf("Invalid authentication code for user %v", username)
	}
	user.Validated = true
	err := r.DB.WithContext(ctx).Save(user).Error
	return err
}

func (r *UserRepo) CreateContact(ctx context.Context, user1 *models.User, user2 *models.User) (*models.Contact, error) {
	contact := &models.Contact{
		FirstUser:  *user1,
		SecondUser: *user2,
		Username1:  user1.Username,
		Username2:  user2.Username,
		CreatedAt:  time.Now(),
	}
	contact_id := r.DB.Create(&contact)
	if contact_id == nil {
		return nil, fmt.Errorf("unable to create a new conact %w for user %w", user2.Username, user1.Username)
	}
	return contact, nil
}
