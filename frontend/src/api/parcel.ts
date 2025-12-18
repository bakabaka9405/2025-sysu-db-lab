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

    // 转换后端字段名到前端格式
    if (response.code === 0 && response.data) {
      response.data.list = response.data.list.map((item: any) => ({
        ...item,
        recipientName: item.recipient_name,
        recipientPhone: item.recipient_phone,
        trackingNumber: item.tracking_number,
        pickupCode: item.pickup_code,
        courierCompany: item.courier_company,
        receivedAt: item.received_at,
        shelvedAt: item.shelved_at,
        pickedUpAt: item.picked_up_at,
        expectedOverdueAt: item.expected_overdue_at,
        shelf: item.shelf ? {
          ...item.shelf,
          shelfCode: item.shelf.shelf_code,
          currentCount: item.shelf.current_count
        } : undefined
      }));
    }

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

    // 转换后端字段名到前端格式
    if (response.code === 0 && response.data) {
      response.data.list = response.data.list.map((item: any) => ({
        ...item,
        recipientName: item.recipient_name,
        recipientPhone: item.recipient_phone,
        trackingNumber: item.tracking_number,
        pickupCode: item.pickup_code,
        courierCompany: item.courier_company,
        receivedAt: item.received_at,
        shelvedAt: item.shelved_at,
        pickedUpAt: item.picked_up_at,
        expectedOverdueAt: item.expected_overdue_at,
        shelf: item.shelf ? {
          ...item.shelf,
          shelfCode: item.shelf.shelf_code,
          currentCount: item.shelf.current_count
        } : undefined
      }));
    }

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

    // 转换后端字段名到前端格式
    if (response.code === 0 && response.data) {
      response.data = {
        ...response.data,
        recipientName: response.data.recipient_name,
        recipientPhone: response.data.recipient_phone,
        trackingNumber: response.data.tracking_number,
        pickupCode: response.data.pickup_code,
        courierCompany: response.data.courier_company,
        receivedAt: response.data.received_at,
        shelvedAt: response.data.shelved_at,
        pickedUpAt: response.data.picked_up_at,
        expectedOverdueAt: response.data.expected_overdue_at,
        shelf: response.data.shelf ? {
          ...response.data.shelf,
          shelfCode: response.data.shelf.shelf_code,
          currentCount: response.data.shelf.current_count
        } : undefined
      };
    }

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

    // 转换后端字段名到前端格式
    if (response.code === 0 && response.data) {
      response.data = {
        ...response.data,
        recipientName: response.data.recipient_name,
        recipientPhone: response.data.recipient_phone,
        trackingNumber: response.data.tracking_number,
        pickupCode: response.data.pickup_code,
        courierCompany: response.data.courier_company,
        receivedAt: response.data.received_at,
        shelvedAt: response.data.shelved_at,
        pickedUpAt: response.data.picked_up_at,
        expectedOverdueAt: response.data.expected_overdue_at,
        shelf: response.data.shelf ? {
          ...response.data.shelf,
          shelfCode: response.data.shelf.shelf_code,
          currentCount: response.data.shelf.current_count
        } : undefined
      };
    }

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

    // 转换后端字段名到前端格式
    if (response.code === 0 && response.data) {
      response.data = {
        ...response.data,
        recipientName: response.data.recipient_name,
        recipientPhone: response.data.recipient_phone,
        trackingNumber: response.data.tracking_number,
        pickupCode: response.data.pickup_code,
        courierCompany: response.data.courier_company,
        receivedAt: response.data.received_at,
        shelvedAt: response.data.shelved_at,
        pickedUpAt: response.data.picked_up_at,
        expectedOverdueAt: response.data.expected_overdue_at,
        shelf: response.data.shelf ? {
          ...response.data.shelf,
          shelfCode: response.data.shelf.shelf_code,
          currentCount: response.data.shelf.current_count
        } : undefined
      };
    }

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
 * 包裹取件
 */
export const pickupParcel = async (data: PickupParcelRequest): Promise<ApiResponse<Parcel>> => {
  try {
    const response: any = await request.post('/parcels/pickup', data);

    // 转换后端字段名到前端格式
    if (response.code === 0 && response.data) {
      response.data = {
        ...response.data,
        recipientName: response.data.recipient_name,
        recipientPhone: response.data.recipient_phone,
        trackingNumber: response.data.tracking_number,
        pickupCode: response.data.pickup_code,
        courierCompany: response.data.courier_company,
        receivedAt: response.data.received_at,
        shelvedAt: response.data.shelved_at,
        pickedUpAt: response.data.picked_up_at,
        expectedOverdueAt: response.data.expected_overdue_at
      };
    }

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
          ready_for_pickup: response.data.parcels.ready_for_pickup,
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
