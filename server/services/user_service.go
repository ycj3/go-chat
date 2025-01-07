package services

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ycj3/go-chat/server/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	user, err := models.GetUserByID(s.db, userID)
	if err != nil {
		return nil, err
	}

	// Update the LastActive and IsOnline fields
	user.LastActive = time.Now().Unix()
	user.IsOnline = true
	if err := s.db.Save(user).Error; err != nil {
		logrus.Error("Error updating user status:", err)
		return nil, err
	}
	logrus.Debug("User status updated for userID:", userID)

	return user, nil
}
