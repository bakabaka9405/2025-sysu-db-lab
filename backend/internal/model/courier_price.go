package model

import (
	"time"

	"gorm.io/gorm"
)

// CourierPrice 快递价格配置
type CourierPrice struct {
	ID              int64          `gorm:"primarykey" json:"id"`
	CourierCompany  string         `gorm:"uniqueIndex;size:100;not null" json:"courier_company"` // 快递公司
	FirstWeight     float64        `gorm:"type:decimal(10,2);default:1" json:"first_weight"`     // 首重(kg)
	FirstPrice      float64        `gorm:"type:decimal(10,2);not null" json:"first_price"`       // 首重价格
	AdditionalPrice float64        `gorm:"type:decimal(10,2);not null" json:"additional_price"`  // 续重价格(每kg)
	IsActive        bool           `gorm:"default:true" json:"is_active"`                        // 是否启用
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CourierPrice) TableName() string {
	return "courier_prices"
}

// CalculateFreight 计算运费
func (c *CourierPrice) CalculateFreight(weight float64) float64 {
	if weight <= c.FirstWeight {
		return c.FirstPrice
	}
	// 续重部分向上取整
	additionalWeight := weight - c.FirstWeight
	additionalUnits := int(additionalWeight)
	if additionalWeight > float64(additionalUnits) {
		additionalUnits++
	}
	return c.FirstPrice + float64(additionalUnits)*c.AdditionalPrice
}
