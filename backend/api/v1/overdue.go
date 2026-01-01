package v1

import "time"

// OverdueInfo 滞留件信息
type OverdueInfo struct {
	ID             int64      `json:"id"`
	ParcelID       int64      `json:"parcel_id"`
	TrackingNumber string     `json:"tracking_number"`
	PickupCode     string     `json:"pickup_code"`
	RecipientName  string     `json:"recipient_name"`
	RecipientPhone string     `json:"recipient_phone"`
	CourierCompany string     `json:"courier_company"`
	Size           string     `json:"size"`
	ShelfCode      string     `json:"shelf_code,omitempty"`
	OverdueStartAt time.Time  `json:"overdue_start_at"`
	OverdueDays    int        `json:"overdue_days"`
	ReminderCount  int        `json:"reminder_count"`
	LastReminderAt *time.Time `json:"last_reminder_at,omitempty"`
	OverdueFee     float64    `json:"overdue_fee"`
	Status         string     `json:"status"`
}

// OverdueListRequest 滞留件列表请求
type OverdueListRequest struct {
	Status   string `form:"status" example:"pending"`
	Page     int    `form:"page" example:"1"`
	PageSize int    `form:"page_size" example:"20"`
}

// OverdueListData 滞留件列表数据
type OverdueListData struct {
	List       []OverdueInfo `json:"list"`
	Pagination Pagination    `json:"pagination"`
}

// BatchResult 批量操作结果
type BatchResult struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

// BatchReturnRequest 批量退回请求
type BatchReturnRequest struct {
	IDs []int64 `json:"ids" binding:"required"`
}

// CheckOverdueData 检查滞留件响应数据
type CheckOverdueData struct {
	MarkedCount int `json:"marked_count"`
}

// ========== 响应结构 ==========

// OverdueListResponse 滞留件列表响应
type OverdueListResponse struct {
	Response
	Data OverdueListData `json:"data"`
}

// OverdueDetailResponse 滞留件详情响应
type OverdueDetailResponse struct {
	Response
	Data OverdueInfo `json:"data"`
}

// BatchReturnResponse 批量退回响应
type BatchReturnResponse struct {
	Response
	Data BatchResult `json:"data"`
}

// CheckOverdueResponse 检查滞留件响应
type CheckOverdueResponse struct {
	Response
	Data CheckOverdueData `json:"data"`
}
