package v1

import "time"

// CreateShipmentRequest 创建寄件单请求
type CreateShipmentRequest struct {
	RecipientName    string  `json:"recipient_name" binding:"required" example:"李四"`
	RecipientPhone   string  `json:"recipient_phone" binding:"required" example:"13900139000"`
	RecipientAddress string  `json:"recipient_address" binding:"required" example:"广东省广州市天河区XX路XX号"`
	CourierCompany   string  `json:"courier_company" binding:"required" example:"顺丰速运"`
	Weight           float64 `json:"weight" binding:"required,gt=0" example:"2.5"`
	Volume           float64 `json:"volume" binding:"required,gt=0" example:"0.05"`
}

// CalculateFreightRequest 运费计算请求
type CalculateFreightRequest struct {
	CourierCompany string  `json:"courier_company" binding:"required" example:"顺丰速运"`
	Weight         float64 `json:"weight" binding:"required,gt=0" example:"2.5"`
	Volume         float64 `json:"volume" binding:"required,gt=0" example:"0.05"`
}

// CalculateFreightData 运费计算响应数据
type CalculateFreightData struct {
	Freight             float64 `json:"freight"`
	CalculationDetails  struct {
		BaseFreight   float64 `json:"base_freight"`
		WeightFreight float64 `json:"weight_freight"`
		VolumeFreight float64 `json:"volume_freight"`
	} `json:"calculation_details"`
}

// ShipmentInfo 寄件单信息
type ShipmentInfo struct {
	ID               int64      `json:"id"`
	ShipmentNumber   string     `json:"shipment_number"`
	SenderName       string     `json:"sender_name"`
	SenderPhone      string     `json:"sender_phone"`
	RecipientName    string     `json:"recipient_name"`
	RecipientPhone   string     `json:"recipient_phone"`
	RecipientAddress string     `json:"recipient_address"`
	CourierCompany   string     `json:"courier_company"`
	Weight           float64    `json:"weight"`
	Volume           float64    `json:"volume"`
	Freight          float64    `json:"freight"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	ShippedAt        *time.Time `json:"shipped_at,omitempty"`
}

// ShipmentListRequest 寄件单列表请求
type ShipmentListRequest struct {
	Status   string `form:"status" example:"pending"`
	Page     int    `form:"page" example:"1"`
	PageSize int    `form:"page_size" example:"20"`
}

// ShipmentListData 寄件单列表数据
type ShipmentListData struct {
	List       []ShipmentInfo `json:"list"`
	Pagination Pagination     `json:"pagination"`
}
