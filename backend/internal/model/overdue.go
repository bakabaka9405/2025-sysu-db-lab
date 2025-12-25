package model

import (
	"time"

	"gorm.io/gorm"
)

// OverdueStatus 滞留记录状态
type OverdueStatus string

const (
	OverdueStatusPending  OverdueStatus = "pending"  // 待处理
	OverdueStatusReminded OverdueStatus = "reminded" // 已提醒
	OverdueStatusReturned OverdueStatus = "returned" // 已退回
)

// OverdueRecord 滞留记录
type OverdueRecord struct {
	ID             int64          `gorm:"primarykey" json:"id"`
	ParcelID       int64          `gorm:"uniqueIndex;not null" json:"parcel_id"`
	Parcel         *Parcel        `gorm:"foreignKey:ParcelID" json:"parcel,omitempty"`
	OverdueStartAt time.Time      `gorm:"not null" json:"overdue_start_at"`
	ReminderCount  int            `gorm:"default:0" json:"reminder_count"`
	LastReminderAt *time.Time     `json:"last_reminder_at,omitempty"`
	OverdueFee     float64        `gorm:"type:decimal(10,2);default:0" json:"overdue_fee"`
	Status         OverdueStatus  `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (OverdueRecord) TableName() string {
	return "overdue_records"
}

// CalculateFee 计算滞留费用（每天1元）
func (o *OverdueRecord) CalculateFee() float64 {
	days := int(time.Since(o.OverdueStartAt).Hours() / 24)
	if days < 0 {
		days = 0
	}
	return float64(days)
}
