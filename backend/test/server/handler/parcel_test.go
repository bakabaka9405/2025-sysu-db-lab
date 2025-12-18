package handler

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	v1 "backend/api/v1"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/model"
	mock_service "backend/test/mocks/service"
)

func TestParcelHandler_ReceiveParcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.ReceiveParcelRequest{
		TrackingNumber:  "TEST123456",
		RecipientName:   "Test User",
		RecipientPhone:  "13800138000",
		RecipientIDCard: "440101199001011234",
		CourierCompany:  "SF Express",
		Size:            "small",
		Weight:          1.5,
		Notes:           "Test notes",
	}

	shelfID := int64(1)
	mockParcel := &model.Parcel{
		ID:             1,
		TrackingNumber: params.TrackingNumber,
		PickupCode:     "A12345",
		RecipientName:  params.RecipientName,
		RecipientPhone: params.RecipientPhone,
		Status:         model.ParcelStatusReceived,
		ShelfID:        &shelfID,
		ReceivedAt:     time.Now(),
	}

	mockParcelService := mock_service.NewMockParcelService(ctrl)
	mockParcelService.EXPECT().ReceiveParcel(gomock.Any(), &params).Return(mockParcel, nil)

	parcelHandler := handler.NewParcelHandler(hdl, mockParcelService)
	router.POST("/parcels/receive", parcelHandler.ReceiveParcel)

	e := newHttpExcept(t, router)
	obj := e.POST("/parcels/receive").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("tracking_number").IsEqual(params.TrackingNumber)
	objData.Value("pickup_code").IsEqual("A12345")
}

func TestParcelHandler_PickupParcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.PickupParcelRequest{
		PickupCode:     "A12345",
		RecipientPhone: "13800138000",
	}

	pickupTime := time.Now()
	mockParcel := &model.Parcel{
		ID:             1,
		PickupCode:     params.PickupCode,
		RecipientPhone: params.RecipientPhone,
		Status:         model.ParcelStatusPickedUp,
		PickedUpAt:     &pickupTime,
		ReceivedAt:     time.Now(),
	}

	mockParcelService := mock_service.NewMockParcelService(ctrl)
	mockParcelService.EXPECT().PickupParcel(gomock.Any(), &params).Return(mockParcel, nil)

	parcelHandler := handler.NewParcelHandler(hdl, mockParcelService)
	router.POST("/parcels/pickup", parcelHandler.PickupParcel)

	obj := newHttpExcept(t, router).POST("/parcels/pickup").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("pickup_code").IsEqual(params.PickupCode)
	objData.Value("status").IsEqual("picked_up")
}

func TestParcelHandler_GetParcelByPickupCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pickupCode := "A12345"
	mockParcel := &model.Parcel{
		ID:         1,
		PickupCode: pickupCode,
		Status:     model.ParcelStatusReceived,
		ReceivedAt: time.Now(),
	}

	mockParcelService := mock_service.NewMockParcelService(ctrl)
	mockParcelService.EXPECT().GetParcelByPickupCode(gomock.Any(), pickupCode).Return(mockParcel, nil)

	parcelHandler := handler.NewParcelHandler(hdl, mockParcelService)
	router.GET("/parcels/by-pickup-code/:code", parcelHandler.GetParcelByPickupCode)

	obj := newHttpExcept(t, router).GET("/parcels/by-pickup-code/" + pickupCode).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("pickup_code").IsEqual(pickupCode)
}

func TestParcelHandler_GetParcelByPickupCode_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pickupCode := "INVALID"

	mockParcelService := mock_service.NewMockParcelService(ctrl)
	mockParcelService.EXPECT().GetParcelByPickupCode(gomock.Any(), pickupCode).Return(nil, errors.New("not found"))

	parcelHandler := handler.NewParcelHandler(hdl, mockParcelService)
	testRouter := gin.Default()
	testRouter.GET("/parcels/by-pickup-code/:code", parcelHandler.GetParcelByPickupCode)

	obj := newHttpExcept(t, testRouter).GET("/parcels/by-pickup-code/" + pickupCode).
		Expect().
		Status(http.StatusNotFound).
		JSON().
		Object()
	obj.Value("code").NotEqual(0)
}

func TestParcelHandler_GetParcelByTrackingNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackingNumber := "TEST123456"
	mockParcel := &model.Parcel{
		ID:             1,
		TrackingNumber: trackingNumber,
		Status:         model.ParcelStatusReceived,
		ReceivedAt:     time.Now(),
	}

	mockParcelService := mock_service.NewMockParcelService(ctrl)
	mockParcelService.EXPECT().GetParcelByTrackingNumber(gomock.Any(), trackingNumber).Return(mockParcel, nil)

	parcelHandler := handler.NewParcelHandler(hdl, mockParcelService)
	router.GET("/parcels/by-tracking/:number", parcelHandler.GetParcelByTrackingNumber)

	obj := newHttpExcept(t, router).GET("/parcels/by-tracking/" + trackingNumber).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("tracking_number").IsEqual(trackingNumber)
}

func TestParcelHandler_GetParcelList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockData := &v1.ParcelListData{
		List: []v1.ParcelInfo{
			{
				ID:         1,
				PickupCode: "A12345",
				Status:     "received",
			},
		},
		Pagination: v1.Pagination{
			Page:       1,
			PageSize:   20,
			Total:      1,
			TotalPages: 1,
		},
	}

	mockParcelService := mock_service.NewMockParcelService(ctrl)
	mockParcelService.EXPECT().GetParcelList(gomock.Any(), gomock.Any()).Return(mockData, nil)

	parcelHandler := handler.NewParcelHandler(hdl, mockParcelService)
	router.GET("/parcels", parcelHandler.GetParcelList)

	obj := newHttpExcept(t, router).GET("/parcels").
		WithQuery("page", 1).
		WithQuery("page_size", 20).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("pagination").Object().Value("total").IsEqual(1)
}

func TestParcelHandler_GetMyParcels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().GetProfile(gomock.Any(), int64(123)).Return(&v1.GetProfileResponseData{
		ID:    123,
		Phone: "13800138000",
	}, nil)

	mockData := &v1.ParcelListData{
		List: []v1.ParcelInfo{
			{
				ID:         1,
				PickupCode: "A12345",
				Status:     "received",
			},
		},
		Pagination: v1.Pagination{
			Page:       1,
			PageSize:   20,
			Total:      1,
			TotalPages: 1,
		},
	}

	mockParcelService := mock_service.NewMockParcelService(ctrl)
	mockParcelService.EXPECT().GetMyParcels(gomock.Any(), "13800138000", gomock.Any()).Return(mockData, nil)

	testHandler := handler.NewHandler(logger, mockUserService)
	parcelHandler := handler.NewParcelHandler(testHandler, mockParcelService)
	router.Use(middleware.NoStrictAuth(jwt, logger))
	router.GET("/my-parcels", parcelHandler.GetMyParcels)

	obj := newHttpExcept(t, router).GET("/my-parcels").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithQuery("page", 1).
		WithQuery("page_size", 20).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("pagination").Object().Value("total").IsEqual(1)
}

func TestParcelHandler_GetDashboardStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStats := &v1.DashboardStatistics{}
	mockStats.Today.ReceivedCount = 10
	mockStats.Today.PickedUpCount = 5
	mockStats.Parcels.Total = 100

	mockParcelService := mock_service.NewMockParcelService(ctrl)
	mockParcelService.EXPECT().GetDashboardStatistics(gomock.Any()).Return(mockStats, nil)

	parcelHandler := handler.NewParcelHandler(hdl, mockParcelService)
	router.GET("/dashboard/statistics", parcelHandler.GetDashboardStatistics)

	obj := newHttpExcept(t, router).GET("/dashboard/statistics").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("today").Object().Value("received_count").IsEqual(10)
	objData.Value("today").Object().Value("picked_up_count").IsEqual(5)
	objData.Value("parcels").Object().Value("total").IsEqual(100)
}
