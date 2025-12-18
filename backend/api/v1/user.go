package v1

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20" example:"zhangsan"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
	RealName string `json:"realName" binding:"required" example:"张三"`
	Phone    string `json:"phone" binding:"required" example:"13800138000"`
	IDCard   string `json:"idCard" example:"440101199001011234"`
	Email    string `json:"email" binding:"omitempty,email" example:"zhangsan@example.com"`
}

// LoginRequest 登录请求（支持用户名或手机号登录）
type LoginRequest struct {
	UsernameOrPhone string `json:"usernameOrPhone" binding:"required" example:"zhangsan"`
	Password        string `json:"password" binding:"required" example:"123456"`
}

// LoginResponseData 登录响应数据
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
	User        UserInfo `json:"user"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Response
	Data LoginResponseData
}

// UserInfo 用户基本信息
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	Phone    string `json:"phone"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	RealName string `json:"realName" example:"张三"`
	Email    string `json:"email" binding:"omitempty,email" example:"zhangsan@example.com"`
}

// GetProfileResponseData 获取用户信息响应数据
type GetProfileResponseData struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	Phone    string `json:"phone"`
	IDCard   string `json:"idCard,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role"`
}

// GetProfileResponse 获取用户信息响应
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}
