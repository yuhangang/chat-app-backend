package repository

import (
	"context"
	"example/user/hello/internal/db/tables"

	"gorm.io/gorm"
)

type UserRepo struct {
	conn *gorm.DB
}

func NewUserRepo(conn *gorm.DB) *UserRepo {
	return &UserRepo{conn: conn}
}

func (db *UserRepo) CreateUser(ctx context.Context, user tables.User) (tables.User, error) {
	res := db.conn.WithContext(ctx).Create(&user)

	if res.Error != nil {
		return tables.User{}, res.Error // Fix: Use res.Error
	}

	return user, nil // Fix: Return the user object directly
}

func (db *UserRepo) GetUser(ctx context.Context, userID uint) (tables.User, error) {
	var user tables.User

	err := db.conn.WithContext(ctx).Where("id = ?", userID).First(&user).Error

	return user, err
}
