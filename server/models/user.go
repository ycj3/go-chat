package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	UserID     string `gorm:"primaryKey" json:"user_id"`
	Nickname   string `json:"nickname"`
	ProfileURL string `json:"profile_url"`
}

var ErrUserNotFound = errors.New("user not found")

func GetUserByID(db *gorm.DB, userID string) (*User, error) {
	var user User
	if err := db.First(&user, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create the user if not found
			user = User{
				UserID:     userID,
				Nickname:   userID, // Default nickname to userID
				ProfileURL: "",     // Default profile URL to empty string
			}
			if err := db.Create(&user).Error; err != nil {
				return nil, err
			}
			return &user, nil
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}
