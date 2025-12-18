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

func TestShipmentHandler_CreateShipment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.CreateShipmentRequest{
		RecipientName:    "Test Recipient",
		RecipientPhone:   "13900139000",
		RecipientAddress: "Test Address",
		CourierCompany:   "顺丰速运",
		Weight:           2.0,
		Volume:           0.05,
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().GetProfile(gomock.Any(), int64(123)).Return(&v1.GetProfileResponseData{
		ID:       123,
		RealName: "Test Sender",
		Phone:    "13800138000",
	}, nil)

	mockShipment := &model.Shipment{
		ID:               1,
		ShipmentNumber:   "SH202412345678",
		SenderID:         123,
		SenderName:       "Test Sender",
		SenderPhone:      "13800138000",
		RecipientName:    params.RecipientName,
		RecipientPhone:   params.RecipientPhone,
		RecipientAddress: params.RecipientAddress,
		CourierCompany:   params.CourierCompany,
		Weight:           params.Weight,
		Volume:           params.Volume,
		Freight:          20.0,
		Status:           model.ShipmentStatusPending,
		CreatedAt:        time.Now(),
	}

	mockShipmentService := mock_service.NewMockShipmentService(ctrl)
	mockShipmentService.EXPECT().CreateShipment(gomock.Any(), int64(123), "Test Sender", "13800138000", &params).
		Return(mockShipment, nil)

	testHandler := handler.NewHandler(logger, mockUserService)
	shipmentHandler := handler.NewShipmentHandler(testHandler, mockShipmentService)
	router.Use(middleware.NoStrictAuth(jwt, logger))
	router.POST("/shipments", shipmentHandler.CreateShipment)

	obj := newHttpExcept(t, router).POST("/shipments").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("shipment_number").IsEqual("SH202412345678")
}

func TestShipmentHandler_GetShipmentDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipment := &model.Shipment{
		ID:             1,
		ShipmentNumber: "SH202412345678",
		SenderID:       123,
		Status:         model.ShipmentStatusPending,
		CreatedAt:      time.Now(),
	}

	mockShipmentService := mock_service.NewMockShipmentService(ctrl)
	mockShipmentService.EXPECT().GetShipmentByID(gomock.Any(), int64(1)).Return(mockShipment, nil)

	shipmentHandler := handler.NewShipmentHandler(hdl, mockShipmentService)
	router.Use(middleware.NoStrictAuth(jwt, logger))
	router.GET("/shipments/:id", shipmentHandler.GetShipmentDetail)

	obj := newHttpExcept(t, router).GET("/shipments/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("shipment_number").IsEqual("SH202412345678")
}

func TestShipmentHandler_GetShipmentDetail_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipmentService := mock_service.NewMockShipmentService(ctrl)
	mockShipmentService.EXPECT().GetShipmentByID(gomock.Any(), int64(999)).Return(nil, errors.New("not found"))

	shipmentHandler := handler.NewShipmentHandler(hdl, mockShipmentService)
	testRouter := gin.Default()
	testRouter.Use(middleware.NoStrictAuth(jwt, logger))
	testRouter.GET("/shipments/:id", shipmentHandler.GetShipmentDetail)

	obj := newHttpExcept(t, testRouter).GET("/shipments/999").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusNotFound).
		JSON().
		Object()
	obj.Value("code").NotEqual(0)
}

func TestShipmentHandler_GetMyShipments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockData := &v1.ShipmentListData{
		List: []v1.ShipmentInfo{
			{
				ID:             1,
				ShipmentNumber: "SH202412345678",
				Status:         "pending",
			},
		},
		Pagination: v1.Pagination{
			Page:       1,
			PageSize:   20,
			Total:      1,
			TotalPages: 1,
		},
	}

	mockShipmentService := mock_service.NewMockShipmentService(ctrl)
	mockShipmentService.EXPECT().GetMyShipments(gomock.Any(), int64(123), gomock.Any()).Return(mockData, nil)

	shipmentHandler := handler.NewShipmentHandler(hdl, mockShipmentService)
	router.Use(middleware.NoStrictAuth(jwt, logger))
	router.GET("/my-shipments", shipmentHandler.GetMyShipments)

	obj := newHttpExcept(t, router).GET("/my-shipments").
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

func TestShipmentHandler_CalculateFreight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.CalculateFreightRequest{
		CourierCompany: "顺丰速运",
		Weight:         2.0,
		Volume:         0.05,
	}

	mockData := &v1.CalculateFreightData{
		Freight: 20.0,
	}
	mockData.CalculationDetails.BaseFreight = 15.0
	mockData.CalculationDetails.WeightFreight = 5.0
	mockData.CalculationDetails.VolumeFreight = 0.0

	mockShipmentService := mock_service.NewMockShipmentService(ctrl)
	mockShipmentService.EXPECT().CalculateFreight(gomock.Any(), &params).Return(mockData, nil)

	shipmentHandler := handler.NewShipmentHandler(hdl, mockShipmentService)
	router.POST("/shipments/calculate-freight", shipmentHandler.CalculateFreight)

	obj := newHttpExcept(t, router).POST("/shipments/calculate-freight").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("freight").IsEqual(20.0)
}

func TestShipmentHandler_GetCourierCompanies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompanies := []string{"顺丰速运", "中通快递", "圆通速递"}

	mockShipmentService := mock_service.NewMockShipmentService(ctrl)
	mockShipmentService.EXPECT().GetCourierCompanies(gomock.Any()).Return(mockCompanies)

	shipmentHandler := handler.NewShipmentHandler(hdl, mockShipmentService)
	router.GET("/courier-companies", shipmentHandler.GetCourierCompanies)

	obj := newHttpExcept(t, router).GET("/courier-companies").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	obj.Value("data").Array().Length().IsEqual(3)
}
