package model

import (
	"time"

	"gorm.io/gorm"
)

// ParcelStatus 包裹状态
type ParcelStatus string

const (
	ParcelStatusReceived      ParcelStatus = "received"        // 已接收
	ParcelStatusShelved       ParcelStatus = "shelved"         // 已上架
	ParcelStatusReadyForPickup ParcelStatus = "ready_for_pickup" // 待取件
	ParcelStatusPickedUp      ParcelStatus = "picked_up"       // 已取件
	ParcelStatusOverdue       ParcelStatus = "overdue"         // 滞留
	ParcelStatusReturned      ParcelStatus = "returned"        // 已退回
)

// ParcelSize 包裹尺寸
type ParcelSize string

const (
	ParcelSizeSmall  ParcelSize = "small"
	ParcelSizeMedium ParcelSize = "medium"
	ParcelSizeLarge  ParcelSize = "large"
)

// Parcel 包裹信息
type Parcel struct {
	ID               int64        `gorm:"primarykey" json:"id"`
	TrackingNumber   string       `gorm:"uniqueIndex;size:100;not null" json:"tracking_number"`   // 快递单号
	PickupCode       string       `gorm:"uniqueIndex;size:20;not null" json:"pickup_code"`        // 取件码
	RecipientName    string       `gorm:"size:100;not null" json:"recipient_name"`                // 收件人姓名
	RecipientPhone   string       `gorm:"index;size:20;not null" json:"recipient_phone"`          // 收件人手机号
	RecipientIDCard  string       `gorm:"size:20" json:"recipient_id_card,omitempty"`             // 收件人身份证号
	CourierCompany   string       `gorm:"size:50;not null" json:"courier_company"`                // 快递公司
	Size             ParcelSize   `gorm:"type:varchar(20);not null" json:"size"`                  // 包裹尺寸
	Weight           float64      `gorm:"type:decimal(10,2);not null" json:"weight"`              // 重量（kg）
	Status           ParcelStatus `gorm:"type:varchar(20);index;not null" json:"status"`          // 状态
	ShelfID          *int64       `gorm:"index" json:"shelf_id,omitempty"`                        // 货架ID
	Shelf            *Shelf       `gorm:"foreignKey:ShelfID" json:"shelf,omitempty"`              // 货架信息
	ReceivedAt       time.Time    `gorm:"index;not null" json:"received_at"`                      // 入库时间
	ShelvedAt        *time.Time   `json:"shelved_at,omitempty"`                                   // 上架时间
	PickedUpAt       *time.Time   `json:"picked_up_at,omitempty"`                                 // 取件时间
	ExpectedOverdueAt *time.Time  `json:"expected_overdue_at,omitempty"`                          // 预计滞留时间
	Notes            string       `gorm:"type:text" json:"notes,omitempty"`                       // 备注
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Parcel) TableName() string {
	return "parcels"
}
