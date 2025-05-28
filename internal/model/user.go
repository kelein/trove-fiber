package model

import (
	"time"

	"gorm.io/gorm"
)

// User mapped from table <users>
type User struct {
	ID        int32          `gorm:"column:id;primaryKey" json:"id"`
	UserID    string         `gorm:"column:user_id;not null" json:"user_id"`
	Nickname  string         `gorm:"column:nickname;not null" json:"nickname"`
	Password  string         `gorm:"column:password;not null" json:"password"`
	Email     string         `gorm:"column:email;not null" json:"email"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName returns the table name for the User model
func (u *User) TableName() string {
	return "users"
}
