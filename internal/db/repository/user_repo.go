package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/yuhangang/chat-app-backend/internal/db/tables"
	api_errors "github.com/yuhangang/chat-app-backend/user_errors"

	"gorm.io/gorm"
)

type UserRepo struct {
	conn *gorm.DB
}

func NewUserRepo(conn *gorm.DB) *UserRepo {
	return &UserRepo{conn: conn}
}

func (repo *UserRepo) CreateUser(ctx context.Context, user tables.User) (tables.User, error) {
	res := repo.conn.WithContext(ctx).Create(&user)

	if res.Error != nil {
		return tables.User{}, repo.handlUserRepoError(res.Error)
	}

	return user, nil // Fix: Return the user object directly
}

func (repo *UserRepo) GetUser(ctx context.Context, userID uint) (tables.User, error) {
	var user tables.User

	err := repo.conn.WithContext(ctx).Where("id = ?", userID).First(&user).Error

	return user, repo.handlUserRepoError(err)
}

func (repo *UserRepo) GetUserByUsername(ctx context.Context, username string) (tables.User, error) {
	var user tables.User

	err := repo.conn.WithContext(ctx).Where("username = ?", username).First(&user).Error

	return user, repo.handlUserRepoError(err)
}

func (repo *UserRepo) BindUser(ctx context.Context, userID uint, username string) (tables.User, error) {
	var user tables.User

	err := repo.conn.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return tables.User{}, repo.handlUserRepoError(err)
	}

	err = repo.conn.WithContext(ctx).Model(&user).Update("username", username).Error
	if err != nil {

		return tables.User{}, repo.handlUserRepoError(err)
	}

	return user, nil
}

func (repo *UserRepo) handlUserRepoError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return api_errors.ErrUserNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return api_errors.ErrUsernameExists
	}
	return fmt.Errorf("internal error: %w", err)
}
