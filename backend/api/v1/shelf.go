package v1

// CreateShelfRequest 创建货架请求
type CreateShelfRequest struct {
	ShelfCode string `json:"shelf_code" binding:"required" example:"A-1-1"`
	Area      string `json:"area" binding:"required" example:"A"`
	Floor     int    `json:"floor" binding:"required,min=1" example:"1"`
	Column    int    `json:"column" binding:"required,min=1" example:"1"`
	Type      string `json:"type" binding:"required,oneof=small medium large" example:"small"`
	Capacity  int    `json:"capacity" binding:"required,min=1" example:"20"`
}

// UpdateShelfRequest 更新货架请求
type UpdateShelfRequest struct {
	Area     string `json:"area" example:"A"`
	Floor    int    `json:"floor" example:"1"`
	Column   int    `json:"column" example:"1"`
	Type     string `json:"type" example:"small"`
	Capacity int    `json:"capacity" example:"20"`
	Status   string `json:"status" example:"active"`
}

// ShelfListRequest 货架列表请求
type ShelfListRequest struct {
	Area     string `form:"area" example:"A"`
	Type     string `form:"type" example:"small"`
	Status   string `form:"status" example:"active"`
	Page     int    `form:"page" example:"1"`
	PageSize int    `form:"page_size" example:"20"`
}

// ShelfListData 货架列表数据
type ShelfListData struct {
	List       []ShelfInfo `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

// ShelfStatistics 货架统计数据
type ShelfStatistics struct {
	TotalShelves    int     `json:"total_shelves"`
	ActiveShelves   int     `json:"active_shelves"`
	TotalCapacity   int     `json:"total_capacity"`
	TotalUsed       int     `json:"total_used"`
	UtilizationRate float64 `json:"utilization_rate"`
	ByType          []struct {
		Type            string  `json:"type"`
		Count           int     `json:"count"`
		Capacity        int     `json:"capacity"`
		Used            int     `json:"used"`
		UtilizationRate float64 `json:"utilization_rate"`
	} `json:"by_type"`
}
