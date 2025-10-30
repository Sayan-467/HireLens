package models

import "time"

type User struct {
	Id           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	Password string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
