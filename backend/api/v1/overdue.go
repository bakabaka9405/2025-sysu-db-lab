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
