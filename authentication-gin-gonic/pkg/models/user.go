package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Phone     string    `gorm:"uniqueIndex;not null" json:"phone"`
	UserType  string    `gorm:"type:varchar(20);default:'user'" json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
