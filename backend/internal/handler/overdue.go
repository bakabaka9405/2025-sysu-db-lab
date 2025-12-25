package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	v1 "backend/api/v1"
	"backend/internal/service"
)

type OverdueHandler struct {
	*Handler
	overdueService service.OverdueService
}

func NewOverdueHandler(handler *Handler, overdueService service.OverdueService) *OverdueHandler {
	return &OverdueHandler{
		Handler:        handler,
		overdueService: overdueService,
	}
}

// GetOverdueList godoc
// @Summary 获取滞留件列表
// @Schemes
// @Description 获取滞留件列表，支持按状态筛选
// @Tags 滞留件管理
// @Accept json
// @Produce json
// @Param status query string false "状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} v1.Response{data=v1.OverdueListData}
// @Router /overdue [get]
func (h *OverdueHandler) GetOverdueList(ctx *gin.Context) {
	var req v1.OverdueListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.overdueService.List(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// SendReminder godoc
// @Summary 发送滞留提醒
// @Schemes
// @Description 对指定滞留件发送提醒通知
// @Tags 滞留件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "滞留记录ID"
// @Success 200 {object} v1.Response
// @Router /overdue/{id}/remind [post]
func (h *OverdueHandler) SendReminder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.overdueService.SendReminder(ctx, id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// BatchReturn godoc
// @Summary 批量退回滞留件
// @Schemes
// @Description 批量将滞留件退回
// @Tags 滞留件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.BatchReturnRequest true "退回请求"
// @Success 200 {object} v1.Response{data=v1.BatchResult}
// @Router /overdue/batch-return [post]
func (h *OverdueHandler) BatchReturn(ctx *gin.Context) {
	var req v1.BatchReturnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	result, err := h.overdueService.BatchReturn(ctx, req.IDs)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, result)
}

// CheckOverdue godoc
// @Summary 检查并标记滞留件
// @Schemes
// @Description 手动触发检查滞留件（通常由定时任务调用）
// @Tags 滞留件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.Response
// @Router /overdue/check [post]
func (h *OverdueHandler) CheckOverdue(ctx *gin.Context) {
	count, err := h.overdueService.CheckAndMarkOverdue(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, map[string]int{"marked_count": count})
}
