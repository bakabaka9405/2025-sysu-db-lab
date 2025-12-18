package model

import (
	"time"

	"gorm.io/gorm"
)

// ShipmentStatus 寄件单状态
type ShipmentStatus string

const (
	ShipmentStatusPending   ShipmentStatus = "pending"   // 待发货
	ShipmentStatusShipped   ShipmentStatus = "shipped"   // 已发货
	ShipmentStatusDelivered ShipmentStatus = "delivered" // 已送达
	ShipmentStatusCancelled ShipmentStatus = "cancelled" // 已取消
)

// Shipment 寄件单信息
type Shipment struct {
	ID               int64          `gorm:"primarykey" json:"id"`
	ShipmentNumber   string         `gorm:"uniqueIndex;size:100;not null" json:"shipment_number"` // 寄件单号
	SenderID         int64          `gorm:"index;not null" json:"sender_id"`                      // 寄件人ID
	SenderName       string         `gorm:"size:100;not null" json:"sender_name"`                 // 寄件人姓名
	SenderPhone      string         `gorm:"size:20;not null" json:"sender_phone"`                 // 寄件人手机号
	RecipientName    string         `gorm:"size:100;not null" json:"recipient_name"`              // 收件人姓名
	RecipientPhone   string         `gorm:"size:20;not null" json:"recipient_phone"`              // 收件人手机号
	RecipientAddress string         `gorm:"size:500;not null" json:"recipient_address"`           // 收件人地址
	CourierCompany   string         `gorm:"size:50;not null" json:"courier_company"`              // 快递公司
	Weight           float64        `gorm:"type:decimal(10,2);not null" json:"weight"`            // 重量（kg）
	Volume           float64        `gorm:"type:decimal(10,4);not null" json:"volume"`            // 体积（m³）
	Freight          float64        `gorm:"type:decimal(10,2);not null" json:"freight"`           // 运费
	Status           ShipmentStatus `gorm:"type:varchar(20);index;not null" json:"status"`        // 状态
	ShippedAt        *time.Time     `json:"shipped_at,omitempty"`                                 // 发货时间
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Shipment) TableName() string {
	return "shipments"
}
