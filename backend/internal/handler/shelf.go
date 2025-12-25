package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	v1 "backend/api/v1"
	"backend/internal/service"
)

type ShelfHandler struct {
	*Handler
	shelfService service.ShelfService
}

func NewShelfHandler(handler *Handler, shelfService service.ShelfService) *ShelfHandler {
	return &ShelfHandler{
		Handler:      handler,
		shelfService: shelfService,
	}
}

// CreateShelf godoc
// @Summary 创建货架
// @Schemes
// @Description 管理员创建新货架
// @Tags 货架管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.CreateShelfRequest true "货架信息"
// @Success 200 {object} v1.Response{data=v1.ShelfInfo}
// @Router /shelves [post]
func (h *ShelfHandler) CreateShelf(ctx *gin.Context) {
	var req v1.CreateShelfRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	shelf, err := h.shelfService.Create(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, service.ConvertShelfToInfo(shelf))
}

// UpdateShelf godoc
// @Summary 更新货架
// @Schemes
// @Description 管理员更新货架信息
// @Tags 货架管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "货架ID"
// @Param request body v1.UpdateShelfRequest true "货架信息"
// @Success 200 {object} v1.Response{data=v1.ShelfInfo}
// @Router /shelves/{id} [put]
func (h *ShelfHandler) UpdateShelf(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	var req v1.UpdateShelfRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	shelf, err := h.shelfService.Update(ctx, id, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, service.ConvertShelfToInfo(shelf))
}

// DeleteShelf godoc
// @Summary 删除货架
// @Schemes
// @Description 管理员删除货架
// @Tags 货架管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "货架ID"
// @Success 200 {object} v1.Response
// @Router /shelves/{id} [delete]
func (h *ShelfHandler) DeleteShelf(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.shelfService.Delete(ctx, id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// GetShelf godoc
// @Summary 获取货架详情
// @Schemes
// @Description 获取单个货架详情
// @Tags 货架管理
// @Accept json
// @Produce json
// @Param id path int true "货架ID"
// @Success 200 {object} v1.Response{data=v1.ShelfInfo}
// @Router /shelves/{id} [get]
func (h *ShelfHandler) GetShelf(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	shelf, err := h.shelfService.GetByID(ctx, id)
	if err != nil {
		v1.HandleError(ctx, http.StatusNotFound, v1.ErrNotFound, nil)
		return
	}

	v1.HandleSuccess(ctx, service.ConvertShelfToInfo(shelf))
}

// GetShelfList godoc
// @Summary 获取货架列表
// @Schemes
// @Description 获取货架列表，支持筛选
// @Tags 货架管理
// @Accept json
// @Produce json
// @Param area query string false "区域"
// @Param type query string false "类型"
// @Param status query string false "状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} v1.Response{data=v1.ShelfListData}
// @Router /shelves [get]
func (h *ShelfHandler) GetShelfList(ctx *gin.Context) {
	var req v1.ShelfListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.shelfService.List(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// GetShelfStatistics godoc
// @Summary 获取货架统计
// @Schemes
// @Description 获取货架统计数据
// @Tags 货架管理
// @Accept json
// @Produce json
// @Success 200 {object} v1.Response{data=v1.ShelfStatistics}
// @Router /shelves/statistics [get]
func (h *ShelfHandler) GetShelfStatistics(ctx *gin.Context) {
	stats, err := h.shelfService.GetStatistics(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, stats)
}
