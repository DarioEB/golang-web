package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string          `json:"id" gorm:"type:char(36);not null;primary_key;unique_index"`
	Firstname string          `json:"firstname" gorm:"type:char(50);not null"`
	Lastname  string          `json:"lastname" gorm:"type:char(50);not null"`
	Email     string          `json:"email" gorm:"type:char(50);not null"`
	Phone     string          `json:"phone" gorm:"type:char(30);not null"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
