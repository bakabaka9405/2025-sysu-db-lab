import request from './request';
import type {
  Parcel,
  ParcelListResponse,
  ReceiveParcelRequest,
  PickupParcelRequest,
  DashboardStatistics,
  ApiResponse
} from './types';

/**
 * 获取包裹列表(管理端)
 */
export const getParcelList = async (params?: {
  status?: string
  courier_company?: string
  recipient_phone?: string
  page?: number
  page_size?: number
}): Promise<ApiResponse<ParcelListResponse>> => {
  try {
    const response: any = await request.get('/parcels', { params });
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '获取包裹列表失败',
      data: undefined
    };
  }
};

/**
 * 获取我的包裹
 */
export const getMyParcels = async (params?: {
  status?: string
  page?: number
  page_size?: number
}): Promise<ApiResponse<ParcelListResponse>> => {
  try {
    const response: any = await request.get('/my-parcels', { params });
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '获取包裹列表失败',
      data: undefined
    };
  }
};

/**
 * 通过取件码查询包裹
 */
export const getParcelByPickupCode = async (pickupCode: string): Promise<ApiResponse<Parcel>> => {
  try {
    const response: any = await request.get(`/parcels/by-pickup-code/${pickupCode}`);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 40002,
      message: error.response?.data?.message || '包裹不存在',
      data: undefined
    };
  }
};

/**
 * 通过快递单号查询包裹
 */
export const getParcelByTrackingNumber = async (
  trackingNumber: string
): Promise<ApiResponse<Parcel>> => {
  try {
    const response: any = await request.get(`/parcels/by-tracking/${trackingNumber}`);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 40002,
      message: error.response?.data?.message || '包裹不存在',
      data: undefined
    };
  }
};

/**
 * 包裹入库
 */
export const receiveParcel = async (
  data: ReceiveParcelRequest
): Promise<ApiResponse<Parcel>> => {
  try {
    const response: any = await request.post('/parcels/receive', data);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '入库失败',
      data: undefined
    };
  }
};

/**
 * 包裹上架
 */
export const shelveParcel = async (id: number): Promise<ApiResponse<Parcel>> => {
  try {
    const response: any = await request.post(`/parcels/${id}/shelve`);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '上架失败',
      data: undefined
    };
  }
};

/**
 * 包裹取件
 */
export const pickupParcel = async (data: PickupParcelRequest): Promise<ApiResponse<Parcel>> => {
  try {
    const response: any = await request.post('/parcels/pickup', data);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '取件失败',
      data: undefined
    };
  }
};

/**
 * 获取仪表盘统计数据
 */
export const getDashboardStatistics = async (): Promise<
  ApiResponse<DashboardStatistics>
> => {
  try {
    const response: any = await request.get('/dashboard/statistics');

    // 转换后端字段名到前端格式
    if (response.code === 0 && response.data) {
      response.data = {
        today: {
          received_count: response.data.today.received_count,
          picked_up_count: response.data.today.picked_up_count,
          overdue_count: response.data.today.overdue_count
        },
        parcels: {
          total: response.data.parcels.total,
          received: response.data.parcels.received,
          shelved: response.data.parcels.shelved,
          ready: response.data.parcels.ready,
          picked_up: response.data.parcels.picked_up,
          overdue: response.data.parcels.overdue,
          returned: response.data.parcels.returned
        },
        shelves: {
          total_capacity: response.data.shelves.total_capacity,
          total_used: response.data.shelves.total_used,
          utilization_rate: response.data.shelves.utilization_rate
        },
        recent_trends: response.data.recent_trends.map((trend: any) => ({
          date: trend.date,
          received: trend.received,
          picked_up: trend.picked_up
        }))
      };
    }

    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '获取统计数据失败',
      data: undefined
    };
  }
};
