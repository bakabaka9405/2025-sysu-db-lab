package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "backend/api/v1"
	"backend/internal/service"
)

type ShipmentHandler struct {
	*Handler
	shipmentService service.ShipmentService
}

func NewShipmentHandler(handler *Handler, shipmentService service.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{
		Handler:         handler,
		shipmentService: shipmentService,
	}
}

// CreateShipment godoc
// @Summary 创建寄件单
// @Schemes
// @Description 用户创建寄件单
// @Tags 寄件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.CreateShipmentRequest true "寄件信息"
// @Success 200 {object} v1.Response{data=v1.ShipmentInfo}
// @Router /shipments [post]
func (h *ShipmentHandler) CreateShipment(ctx *gin.Context) {
	var req v1.CreateShipmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	userId := GetUserIdFromCtx(ctx)
	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	shipment, err := h.shipmentService.CreateShipment(ctx, userId, user.RealName, user.Phone, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, service.ConvertShipmentToInfo(shipment))
}

// GetShipmentDetail godoc
// @Summary 获取寄件单详情
// @Schemes
// @Description 获取寄件单详细信息
// @Tags 寄件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "寄件单ID"
// @Success 200 {object} v1.Response{data=v1.ShipmentInfo}
// @Router /shipments/{id} [get]
func (h *ShipmentHandler) GetShipmentDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	shipment, err := h.shipmentService.GetShipmentByID(ctx, id)
	if err != nil {
		v1.HandleError(ctx, http.StatusNotFound, v1.ErrNotFound, nil)
		return
	}

	v1.HandleSuccess(ctx, service.ConvertShipmentToInfo(shipment))
}

// GetMyShipments godoc
// @Summary 获取我的寄件单列表
// @Schemes
// @Description 用户获取自己的寄件单列表
// @Tags 寄件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query string false "寄件单状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} v1.Response{data=v1.ShipmentListData}
// @Router /my-shipments [get]
func (h *ShipmentHandler) GetMyShipments(ctx *gin.Context) {
	var req v1.ShipmentListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	userId := GetUserIdFromCtx(ctx)

	data, err := h.shipmentService.GetMyShipments(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// CalculateFreight godoc
// @Summary 运费试算
// @Schemes
// @Description 计算寄件运费
// @Tags 寄件管理
// @Accept json
// @Produce json
// @Param request body v1.CalculateFreightRequest true "运费计算参数"
// @Success 200 {object} v1.Response{data=v1.CalculateFreightData}
// @Router /shipments/calculate-freight [post]
func (h *ShipmentHandler) CalculateFreight(ctx *gin.Context) {
	var req v1.CalculateFreightRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.shipmentService.CalculateFreight(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// GetCourierCompanies godoc
// @Summary 获取快递公司列表
// @Schemes
// @Description 获取支持的快递公司列表
// @Tags 寄件管理
// @Accept json
// @Produce json
// @Success 200 {object} v1.Response{data=[]string}
// @Router /courier-companies [get]
func (h *ShipmentHandler) GetCourierCompanies(ctx *gin.Context) {
	companies := h.shipmentService.GetCourierCompanies(ctx)
	v1.HandleSuccess(ctx, companies)
}
