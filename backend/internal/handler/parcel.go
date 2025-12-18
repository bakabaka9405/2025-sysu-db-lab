package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "backend/api/v1"
	"backend/internal/service"
)

type ParcelHandler struct {
	*Handler
	parcelService service.ParcelService
}

func NewParcelHandler(handler *Handler, parcelService service.ParcelService) *ParcelHandler {
	return &ParcelHandler{
		Handler:       handler,
		parcelService: parcelService,
	}
}

// ReceiveParcel godoc
// @Summary 包裹入库
// @Schemes
// @Description 工作人员将包裹录入系统并上架
// @Tags 包裹管理
// @Accept json
// @Produce json
// @Param request body v1.ReceiveParcelRequest true "入库信息"
// @Success 200 {object} v1.Response{data=v1.ParcelInfo}
// @Router /parcels/receive [post]
func (h *ParcelHandler) ReceiveParcel(ctx *gin.Context) {
	var req v1.ReceiveParcelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	parcel, err := h.parcelService.ReceiveParcel(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, service.ConvertParcelToInfo(parcel))
}

// PickupParcel godoc
// @Summary 包裹取件
// @Schemes
// @Description 用户凭取件码和手机号取件
// @Tags 包裹管理
// @Accept json
// @Produce json
// @Param request body v1.PickupParcelRequest true "取件信息"
// @Success 200 {object} v1.Response{data=v1.ParcelInfo}
// @Router /parcels/pickup [post]
func (h *ParcelHandler) PickupParcel(ctx *gin.Context) {
	var req v1.PickupParcelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	parcel, err := h.parcelService.PickupParcel(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, service.ConvertParcelToInfo(parcel))
}

// GetParcelByPickupCode godoc
// @Summary 通过取件码查询包裹
// @Schemes
// @Description 根据取件码查询包裹信息
// @Tags 包裹管理
// @Accept json
// @Produce json
// @Param code path string true "取件码"
// @Success 200 {object} v1.Response{data=v1.ParcelInfo}
// @Router /parcels/by-pickup-code/{code} [get]
func (h *ParcelHandler) GetParcelByPickupCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	parcel, err := h.parcelService.GetParcelByPickupCode(ctx, code)
	if err != nil {
		v1.HandleError(ctx, http.StatusNotFound, v1.ErrNotFound, nil)
		return
	}

	v1.HandleSuccess(ctx, service.ConvertParcelToInfo(parcel))
}

// GetParcelByTrackingNumber godoc
// @Summary 通过快递单号查询包裹
// @Schemes
// @Description 根据快递单号查询包裹信息
// @Tags 包裹管理
// @Accept json
// @Produce json
// @Param number path string true "快递单号"
// @Success 200 {object} v1.Response{data=v1.ParcelInfo}
// @Router /parcels/by-tracking/{number} [get]
func (h *ParcelHandler) GetParcelByTrackingNumber(ctx *gin.Context) {
	number := ctx.Param("number")
	if number == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	parcel, err := h.parcelService.GetParcelByTrackingNumber(ctx, number)
	if err != nil {
		v1.HandleError(ctx, http.StatusNotFound, v1.ErrNotFound, nil)
		return
	}

	v1.HandleSuccess(ctx, service.ConvertParcelToInfo(parcel))
}

// GetParcelList godoc
// @Summary 获取包裹列表（管理端）
// @Schemes
// @Description 工作人员获取包裹列表，支持筛选
// @Tags 包裹管理
// @Accept json
// @Produce json
// @Param status query string false "包裹状态"
// @Param courier_company query string false "快递公司"
// @Param recipient_phone query string false "收件人手机号"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} v1.Response{data=v1.ParcelListData}
// @Router /parcels [get]
func (h *ParcelHandler) GetParcelList(ctx *gin.Context) {
	var req v1.ParcelListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.parcelService.GetParcelList(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// GetMyParcels godoc
// @Summary 获取我的包裹
// @Schemes
// @Description 用户获取自己的包裹列表
// @Tags 包裹管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query string false "包裹状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} v1.Response{data=v1.ParcelListData}
// @Router /my-parcels [get]
func (h *ParcelHandler) GetMyParcels(ctx *gin.Context) {
	var req v1.ParcelListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	userId := GetUserIdFromCtx(ctx)
	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	data, err := h.parcelService.GetMyParcels(ctx, user.Phone, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// GetDashboardStatistics godoc
// @Summary 获取仪表盘统计数据
// @Schemes
// @Description 获取包裹和货架的统计信息
// @Tags 统计
// @Accept json
// @Produce json
// @Success 200 {object} v1.Response{data=v1.DashboardStatistics}
// @Router /dashboard/statistics [get]
func (h *ParcelHandler) GetDashboardStatistics(ctx *gin.Context) {
	stats, err := h.parcelService.GetDashboardStatistics(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, stats)
}
