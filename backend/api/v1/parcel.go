package v1

import "time"

// ReceiveParcelRequest 包裹入库请求
type ReceiveParcelRequest struct {
	TrackingNumber  string  `json:"tracking_number" binding:"required" example:"SF1234567890"`
	RecipientName   string  `json:"recipient_name" binding:"required" example:"张三"`
	RecipientPhone  string  `json:"recipient_phone" binding:"required" example:"13800138000"`
	RecipientIDCard string  `json:"recipient_id_card" example:"440101199001011234"`
	CourierCompany  string  `json:"courier_company" binding:"required" example:"顺丰速运"`
	Size            string  `json:"size" binding:"required,oneof=small medium large" example:"small"`
	Weight          float64 `json:"weight" binding:"required,gt=0" example:"1.5"`
	Notes           string  `json:"notes" example:"易碎物品"`
}

// PickupParcelRequest 包裹取件请求
type PickupParcelRequest struct {
	PickupCode     string `json:"pickup_code" binding:"required" example:"A12345"`
	RecipientPhone string `json:"recipient_phone" binding:"required" example:"13800138000"`
}

// ShelfInfo 货架信息
type ShelfInfo struct {
	ID         int64  `json:"id"`
	ShelfCode  string `json:"shelf_code"`
	Area       string `json:"area"`
	Floor      int    `json:"floor"`
	Column     int    `json:"column"`
	Type       string `json:"type"`
	Capacity   int    `json:"capacity"`
	CurrentCount int  `json:"current_count,omitempty"`
	Status     string `json:"status,omitempty"`
}

// ParcelInfo 包裹信息
type ParcelInfo struct {
	ID                int64      `json:"id"`
	TrackingNumber    string     `json:"tracking_number"`
	PickupCode        string     `json:"pickup_code"`
	RecipientName     string     `json:"recipient_name"`
	RecipientPhone    string     `json:"recipient_phone"`
	CourierCompany    string     `json:"courier_company"`
	Size              string     `json:"size"`
	Weight            float64    `json:"weight"`
	Status            string     `json:"status"`
	Shelf             *ShelfInfo `json:"shelf,omitempty"`
	ReceivedAt        time.Time  `json:"received_at"`
	ShelvedAt         *time.Time `json:"shelved_at,omitempty"`
	PickedUpAt        *time.Time `json:"picked_up_at,omitempty"`
	ExpectedOverdueAt *time.Time `json:"expected_overdue_at,omitempty"`
	Notes             string     `json:"notes,omitempty"`
}

// ParcelListRequest 包裹列表请求
type ParcelListRequest struct {
	Status         string `form:"status" example:"received"`
	CourierCompany string `form:"courier_company" example:"顺丰速运"`
	RecipientPhone string `form:"recipient_phone" example:"13800138000"`
	Page           int    `form:"page" example:"1"`
	PageSize       int    `form:"page_size" example:"20"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ParcelListData 包裹列表数据
type ParcelListData struct {
	List       []ParcelInfo `json:"list"`
	Pagination Pagination   `json:"pagination"`
}

// DashboardStatistics 仪表盘统计数据
type DashboardStatistics struct {
	Today struct {
		ReceivedCount  int `json:"received_count"`
		PickedUpCount  int `json:"picked_up_count"`
		OverdueCount   int `json:"overdue_count"`
	} `json:"today"`
	Parcels struct {
		Total          int `json:"total"`
		Received       int `json:"received"`
		Shelved        int `json:"shelved"`
		ReadyForPickup int `json:"ready_for_pickup"`
		PickedUp       int `json:"picked_up"`
		Overdue        int `json:"overdue"`
		Returned       int `json:"returned"`
	} `json:"parcels"`
	Shelves struct {
		TotalCapacity   int     `json:"total_capacity"`
		TotalUsed       int     `json:"total_used"`
		UtilizationRate float64 `json:"utilization_rate"`
	} `json:"shelves"`
	RecentTrends []struct {
		Date      string `json:"date"`
		Received  int    `json:"received"`
		PickedUp  int    `json:"picked_up"`
	} `json:"recent_trends"`
}

// ========== 响应结构 ==========

// ReceiveParcelResponse 包裹入库响应
type ReceiveParcelResponse struct {
	Response
	Data ParcelInfo `json:"data"`
}

// PickupParcelResponse 包裹取件响应
type PickupParcelResponse struct {
	Response
	Data ParcelInfo `json:"data"`
}

// ParcelDetailResponse 包裹详情响应
type ParcelDetailResponse struct {
	Response
	Data ParcelInfo `json:"data"`
}

// ParcelListResponse 包裹列表响应
type ParcelListResponse struct {
	Response
	Data ParcelListData `json:"data"`
}

// DashboardStatisticsResponse 仪表盘统计响应
type DashboardStatisticsResponse struct {
	Response
	Data DashboardStatistics `json:"data"`
}
