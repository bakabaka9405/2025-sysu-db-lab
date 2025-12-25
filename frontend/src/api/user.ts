import request from './request';
import type {
  RegisterRequest,
  LoginRequest,
  LoginResponse,
  UserProfile,
  UpdateProfileRequest,
  ApiResponse
} from './types';

/**
 * 用户注册
 */
export const register = async (data: RegisterRequest): Promise<ApiResponse> => {
  try {
    return await request.post('/register', data);
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '注册失败',
      data: undefined
    };
  }
};

/**
 * 用户登录
 */
export const login = async (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
  try {
    const response: any = await request.post('/login', data);

    // 保存token和用户信息到localStorage
    if (response.code === 0 && response.data) {
      localStorage.setItem('accessToken', response.data.accessToken);
      localStorage.setItem('userInfo', JSON.stringify(response.data.user));
    }

    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '登录失败',
      data: undefined
    };
  }
};

/**
 * 获取用户信息
 */
export const getProfile = async (): Promise<ApiResponse<UserProfile>> => {
  try {
    return await request.get('/user');
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '获取用户信息失败',
      data: undefined
    };
  }
};

/**
 * 更新用户信息
 */
export const updateProfile = async (data: UpdateProfileRequest): Promise<ApiResponse> => {
  try {
    const response: any = await request.put('/user', data);

    // 更新localStorage中的用户信息
    if (response.code === 0) {
      const userStr = localStorage.getItem('userInfo');
      if (userStr) {
        const userInfo = JSON.parse(userStr);
        userInfo.realName = data.realName;
        if (data.email !== undefined) {
          userInfo.email = data.email;
        }
        localStorage.setItem('userInfo', JSON.stringify(userInfo));
      }
    }

    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '更新失败',
      data: undefined
    };
  }
};
