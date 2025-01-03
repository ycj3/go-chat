package models

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
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
	logrus.Debug("GetUserByID called with userID:", userID)
	var user User
	if err := db.First(&user, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Debug("User not found, creating new user with userID:", userID)
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
				logrus.Error("Error creating new user:", err)
				return nil, err
			}
			logrus.Debug("New user created with userID:", userID)
			return &user, nil
		}
		logrus.Error("Error retrieving user by ID:", err)
		return nil, err
	}
	// Update the LastActive and IsOnline fields
	user.LastActive = time.Now().Unix()
	user.IsOnline = user.IsCurrentlyOnline()
	if err := db.Save(&user).Error; err != nil {
		logrus.Error("Error updating user:", err)
		return nil, err
	}
	logrus.Debug("User retrieved and updated with userID:", userID)
	return &user, nil
}

// GetAllUsers retrieves all users from the database.
func GetAllUsers(db *gorm.DB) ([]User, error) {
	logrus.Debug("GetAllUsers called")
	var users []User
	if err := db.Find(&users).Error; err != nil {
		logrus.Error("Error retrieving all users:", err)
		return nil, err
	}
	logrus.Debug("All users retrieved, count:", len(users))
	return users, nil
}

// GetOnlineUsers retrieves all online users from the database.
func GetOnlineUsers(db *gorm.DB) ([]User, error) {
	logrus.Debug("GetOnlineUsers called")
	var users []User
	if err := db.Find(&users).Error; err != nil {
		logrus.Error("Error retrieving all users:", err)
		return nil, err
	}

	var onlineUsers []User
	for _, user := range users {
		if user.IsCurrentlyOnline() {
			onlineUsers = append(onlineUsers, user)
		}
	}

	if len(onlineUsers) == 0 {
		logrus.Warn("No online users found")
		return nil, errors.New("no online users found")
	}

	logrus.Debug("Online users retrieved, count:", len(onlineUsers))
	return onlineUsers, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	logrus.Debug("CreateUser called with userID:", user.UserID)
	if err := db.Create(user).Error; err != nil {
		logrus.Error("Error creating user:", err)
		return err
	}
	logrus.Debug("User created with userID:", user.UserID)
	return nil
}

func (u *User) IsCurrentlyOnline() bool {
	return time.Now().Unix()-u.LastActive <= 300
}
