import request from './request';
import type {
  ApiResponse,
  OverdueListResponse,
  OverdueInfo,
  BatchResult,
  BatchReturnRequest
} from './types';

// 重新导出类型供外部使用
export type { OverdueInfo, OverdueListResponse, BatchResult };

/**
 * 获取滞留件列表
 */
export const getOverdueList = async (params?: {
  status?: string;
  page?: number;
  page_size?: number;
}): Promise<ApiResponse<OverdueListResponse>> => {
  try {
    const response: any = await request.get('/overdue', { params });
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '获取滞留件列表失败',
      data: undefined
    };
  }
};

/**
 * 发送提醒
 */
export const sendReminder = async (id: number): Promise<ApiResponse<void>> => {
  try {
    const response: any = await request.post(`/overdue/${id}/remind`);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '发送提醒失败',
      data: undefined
    };
  }
};

/**
 * 批量退回
 */
export const batchReturn = async (ids: number[]): Promise<ApiResponse<BatchResult>> => {
  try {
    const body: BatchReturnRequest = { ids };
    const response: any = await request.post('/overdue/batch-return', body);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '批量退回失败',
      data: undefined
    };
  }
};

/**
 * 检查滞留件
 */
export const checkOverdue = async (): Promise<ApiResponse<{ marked_count: number }>> => {
  try {
    const response: any = await request.post('/overdue/check');
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '检查失败',
      data: undefined
    };
  }
};
