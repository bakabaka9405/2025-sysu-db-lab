package model

import (
	"time"

	"gorm.io/gorm"
)

// ShelfStatus 货架状态
type ShelfStatus string

const (
	ShelfStatusActive   ShelfStatus = "active"   // 启用
	ShelfStatusInactive ShelfStatus = "inactive" // 停用
)

// ShelfType 货架类型
type ShelfType string

const (
	ShelfTypeSmall  ShelfType = "small"  // 小件货架
	ShelfTypeMedium ShelfType = "medium" // 中件货架
	ShelfTypeLarge  ShelfType = "large"  // 大件货架
)

// Shelf 货架信息
type Shelf struct {
	ID           int64       `gorm:"primarykey" json:"id"`
	ShelfCode    string      `gorm:"uniqueIndex;size:50;not null" json:"shelf_code"` // 货架编码
	Area         string      `gorm:"size:50;not null" json:"area"`                   // 区域
	Floor        int         `gorm:"not null" json:"floor"`                          // 层
	Column       int         `gorm:"not null" json:"column"`                         // 列
	Type         ShelfType   `gorm:"type:varchar(20);not null" json:"type"`          // 类型
	Capacity     int         `gorm:"not null" json:"capacity"`                       // 容量
	CurrentCount int         `gorm:"default:0" json:"current_count"`                 // 当前存放数量
	Status       ShelfStatus `gorm:"type:varchar(20);not null" json:"status"`        // 状态
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Shelf) TableName() string {
	return "shelves"
}

// IsFull 判断货架是否已满
func (s *Shelf) IsFull() bool {
	return s.CurrentCount >= s.Capacity
}

// IsAvailable 判断货架是否可用
func (s *Shelf) IsAvailable() bool {
	return s.Status == ShelfStatusActive && !s.IsFull()
}
