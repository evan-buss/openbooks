package model

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

// EnsureDefaultUser creates a default user if one does not exist.
func EnsureDefaultUser(db *gorm.DB, ctx context.Context) error {
	var count int64
	tx := db.WithContext(ctx).Model(&User{}).Count(&count)
	if tx.Error != nil {
		return tx.Error
	}

	if count == 0 {
		password, err := bcrypt.GenerateFromPassword([]byte("admin"), 14)
		if err != nil {
			return err
		}

		user := User{
			Username: "admin",
			Password: password,
		}

		tx := db.WithContext(ctx).Create(&user)
		if tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}

func FindUserById(db *gorm.DB, ctx context.Context, id uint) (User, error) {
	var user User
	tx := db.WithContext(ctx).First(&user, id)
	if tx.Error != nil {
		return User{}, tx.Error
	}
	return user, nil
}

func FindUserByUsername(db *gorm.DB, ctx context.Context, username string) (User, error) {
	var user User
	tx := db.WithContext(ctx).First(&user, "username = ?", username)
	if tx.Error != nil {
		return User{}, tx.Error
	}
	return user, nil
}
