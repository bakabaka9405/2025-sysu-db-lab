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
    return await request.post('/register', {
      username: data.username,
      password: data.password,
      real_name: data.realName,
      phone: data.phone,
      email: data.email,
      id_card: data.idCard
    });
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
    const response: any = await request.post('/login', {
      username_or_phone: data.usernameOrPhone,
      password: data.password
    });

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
    const response: any = await request.get('/user');

    // 更新localStorage中的用户信息
    if (response.code === 0 && response.data) {
      const userInfo = {
        id: response.data.id,
        username: response.data.username,
        realName: response.data.real_name,
        phone: response.data.phone,
        email: response.data.email,
        role: response.data.role
      };
      localStorage.setItem('userInfo', JSON.stringify(userInfo));

      // 转换字段名以匹配前端类型
      return {
        code: response.code,
        message: response.message,
        data: {
          id: response.data.id,
          username: response.data.username,
          realName: response.data.real_name,
          phone: response.data.phone,
          email: response.data.email,
          idCard: response.data.id_card,
          role: response.data.role
        }
      };
    }

    return response;
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
    const response: any = await request.put('/user', {
      real_name: data.realName,
      email: data.email
    });

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
