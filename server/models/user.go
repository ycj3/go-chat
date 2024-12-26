package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID     string `gorm:"primaryKey" json:"user_id"`
	Nickname   string `json:"nickname"`
	ProfileURL string `json:"profile_url"`
	LastActive int64  `json:"last_active"`
	CreatedAt  int64  `json:"created_at"`
	IsOnline   bool   `json:"is_online"`
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
				CreatedAt:  time.Now().Unix(),
				LastActive: time.Now().Unix(),
				IsOnline:   true,
			}
			if err := db.Create(&user).Error; err != nil {
				return nil, err
			}
			return &user, nil
		}
		// Update the LastActive and IsOnline fields
		user.LastActive = time.Now().Unix()
		user.IsOnline = true
		if err := db.Save(&user).Error; err != nil {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}
