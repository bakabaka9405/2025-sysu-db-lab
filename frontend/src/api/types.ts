// 用户信息
export interface UserInfo {
  id: number;
  username: string;
  realName: string;
  phone: string;
  email?: string;
  role: string;
}

// 完整用户信息（包含身份证）
export interface UserProfile extends UserInfo {
  idCard?: string;
}

// API响应基础结构
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}

// 注册请求
export interface RegisterRequest {
  username: string;
  password: string;
  realName: string;
  phone: string;
  idCard?: string;
  email?: string;
}

// 登录请求
export interface LoginRequest {
  usernameOrPhone: string;
  password: string;
}

// 登录响应
export interface LoginResponse {
  accessToken: string;
  user: UserInfo;
}

// 更新用户信息请求
export interface UpdateProfileRequest {
  realName: string;
  email?: string;
}

// ========== 包裹相关类型 ==========
export * from './parcel-types'
