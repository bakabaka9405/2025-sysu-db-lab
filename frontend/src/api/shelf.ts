import request from './request';
import type {
  ApiResponse,
  Shelf,
  ShelfListResponse,
  ShelfStatistics,
  CreateShelfRequest,
  UpdateShelfRequest
} from './types';

/**
 * 获取货架列表
 */
export const getShelfList = async (params?: {
  area?: string;
  type?: string;
  status?: string;
  page?: number;
  page_size?: number;
}): Promise<ApiResponse<ShelfListResponse>> => {
  try {
    const response: any = await request.get('/shelves', { params });
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '获取货架列表失败',
      data: undefined
    };
  }
};

/**
 * 获取货架详情
 */
export const getShelfById = async (id: number): Promise<ApiResponse<Shelf>> => {
  try {
    const response: any = await request.get(`/shelves/${id}`);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 40002,
      message: error.response?.data?.message || '货架不存在',
      data: undefined
    };
  }
};

/**
 * 创建货架
 */
export const createShelf = async (data: CreateShelfRequest): Promise<ApiResponse<Shelf>> => {
  try {
    const response: any = await request.post('/shelves', data);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '创建货架失败',
      data: undefined
    };
  }
};

/**
 * 更新货架
 */
export const updateShelf = async (id: number, data: UpdateShelfRequest): Promise<ApiResponse<Shelf>> => {
  try {
    const response: any = await request.put(`/shelves/${id}`, data);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '更新货架失败',
      data: undefined
    };
  }
};

/**
 * 删除货架
 */
export const deleteShelf = async (id: number): Promise<ApiResponse<void>> => {
  try {
    const response: any = await request.delete(`/shelves/${id}`);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '删除货架失败',
      data: undefined
    };
  }
};

/**
 * 获取货架统计
 */
export const getShelfStatistics = async (): Promise<ApiResponse<ShelfStatistics>> => {
  try {
    const response: any = await request.get('/shelves/statistics');
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '获取统计数据失败',
      data: undefined
    };
  }
};
