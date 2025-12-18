package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户表模型
type User struct {
	ID        int64          `gorm:"primaryKey;column:id" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null;column:username" json:"username"`
	Password  string         `gorm:"not null;column:password" json:"-"` // 不在JSON中显示密码
	RealName  string         `gorm:"column:real_name" json:"realName"`
	Phone     string         `gorm:"uniqueIndex;not null;column:phone" json:"phone"`
	IDCard    string         `gorm:"uniqueIndex;column:id_card" json:"idCard,omitempty"`
	Email     string         `gorm:"column:email" json:"email,omitempty"`
	Role      string         `gorm:"not null;default:'user';column:role" json:"role"` // user, staff, admin
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

func (u *User) TableName() string {
	return "users"
}
