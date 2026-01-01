import request from './request';
import type {
  Shipment,
  ShipmentListResponse,
  CreateShipmentRequest,
  CalculateFreightRequest,
  CalculateFreightResponse,
  ApiResponse
} from './types';

/**
 * 获取我的寄件单列表
 */
export const getMyShipments = async (params?: {
  status?: string
  page?: number
  page_size?: number
}): Promise<ApiResponse<ShipmentListResponse>> => {
  try {
    const response: any = await request.get('/my-shipments', { params });
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '获取寄件单列表失败',
      data: undefined
    };
  }
};

/**
 * 创建寄件单
 */
export const createShipment = async (
  data: CreateShipmentRequest
): Promise<ApiResponse<Shipment>> => {
  try {
    const response: any = await request.post('/shipments', data);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '创建失败',
      data: undefined
    };
  }
};

/**
 * 获取寄件单详情
 */
export const getShipmentDetail = async (id: number): Promise<ApiResponse<Shipment>> => {
  try {
    const response: any = await request.get(`/shipments/${id}`);
    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 40002,
      message: error.response?.data?.message || '寄件单不存在',
      data: undefined
    };
  }
};

/**
 * 运费试算
 */
export const calculateFreight = async (
  data: CalculateFreightRequest
): Promise<ApiResponse<CalculateFreightResponse>> => {
  try {
    const response: any = await request.post('/shipments/calculate-freight', data);

    // 转换后端字段名到前端格式
    if (response.code === 0 && response.data) {
      response.data = {
        freight: response.data.freight,
        calculation_details: {
          base_freight: response.data.calculation_details.base_freight,
          weight_freight: response.data.calculation_details.weight_freight,
          volume_freight: response.data.calculation_details.volume_freight
        }
      };
    }

    return response;
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '计算失败',
      data: undefined
    };
  }
};

/**
 * 获取快递公司列表
 */
export const getCourierCompanies = async (): Promise<ApiResponse<string[]>> => {
  try {
    return await request.get('/courier-companies');
  } catch (error: any) {
    return {
      code: error.response?.data?.code || 50000,
      message: error.response?.data?.message || '获取快递公司列表失败',
      data: undefined
    };
  }
};
