package models

type User struct {
	UserID     string `gorm:"primaryKey" json:"user_id"`
	Nickname   string `json:"nickname"`
	ProfileURL string `json:"profile_url"`
	LastActive int64  `json:"last_active"`
	CreatedAt  int64  `json:"created_at"`
	IsOnline   bool   `json:"is_online"`
}
