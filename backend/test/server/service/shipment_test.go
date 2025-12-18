package service_test

import (
	"context"
	"testing"
	"time"

	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/service"
	mock_repository "backend/test/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestShipmentService_CreateShipment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipmentRepo := mock_repository.NewMockShipmentRepository(ctrl)
	shipmentService := service.NewShipmentService(mockShipmentRepo)

	ctx := context.Background()
	senderID := int64(1)
	senderName := "Test Sender"
	senderPhone := "13800138000"
	req := &v1.CreateShipmentRequest{
		RecipientName:    "Test Recipient",
		RecipientPhone:   "13900139000",
		RecipientAddress: "Test Address",
		CourierCompany:   "顺丰速运",
		Weight:           2.0,
		Volume:           0.05,
	}

	mockShipmentRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, shipment *model.Shipment) error {
			shipment.ID = 1
			return nil
		},
	)

	shipment, err := shipmentService.CreateShipment(ctx, senderID, senderName, senderPhone, req)

	assert.NoError(t, err)
	assert.NotNil(t, shipment)
	assert.Equal(t, senderID, shipment.SenderID)
	assert.Equal(t, senderName, shipment.SenderName)
	assert.Equal(t, senderPhone, shipment.SenderPhone)
	assert.Equal(t, req.RecipientName, shipment.RecipientName)
	assert.Equal(t, req.RecipientPhone, shipment.RecipientPhone)
	assert.Equal(t, req.RecipientAddress, shipment.RecipientAddress)
	assert.Equal(t, req.CourierCompany, shipment.CourierCompany)
	assert.Equal(t, req.Weight, shipment.Weight)
	assert.Equal(t, req.Volume, shipment.Volume)
	assert.Equal(t, model.ShipmentStatusPending, shipment.Status)
	assert.Greater(t, shipment.Freight, 0.0)
}

func TestShipmentService_GetShipmentByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipmentRepo := mock_repository.NewMockShipmentRepository(ctrl)
	shipmentService := service.NewShipmentService(mockShipmentRepo)

	ctx := context.Background()
	shipmentID := int64(1)
	mockShipment := &model.Shipment{
		ID:             shipmentID,
		ShipmentNumber: "SH202412345678",
		Status:         model.ShipmentStatusPending,
		CreatedAt:      time.Now(),
	}

	mockShipmentRepo.EXPECT().GetByID(ctx, shipmentID).Return(mockShipment, nil)

	shipment, err := shipmentService.GetShipmentByID(ctx, shipmentID)

	assert.NoError(t, err)
	assert.NotNil(t, shipment)
	assert.Equal(t, shipmentID, shipment.ID)
	assert.Equal(t, "SH202412345678", shipment.ShipmentNumber)
}

func TestShipmentService_GetMyShipments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipmentRepo := mock_repository.NewMockShipmentRepository(ctrl)
	shipmentService := service.NewShipmentService(mockShipmentRepo)

	ctx := context.Background()
	senderID := int64(1)
	req := &v1.ShipmentListRequest{
		Page:     1,
		PageSize: 20,
		Status:   "pending",
	}

	mockShipments := []*model.Shipment{
		{
			ID:             1,
			ShipmentNumber: "SH202412345678",
			SenderID:       senderID,
			Status:         model.ShipmentStatusPending,
			CreatedAt:      time.Now(),
		},
	}

	mockShipmentRepo.EXPECT().ListBySenderID(ctx, senderID, req.Status, req.Page, req.PageSize).
		Return(mockShipments, int64(1), nil)

	data, err := shipmentService.GetMyShipments(ctx, senderID, req)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 1, len(data.List))
	assert.Equal(t, 1, data.Pagination.Total)
	assert.Equal(t, 1, data.Pagination.Page)
	assert.Equal(t, 20, data.Pagination.PageSize)
}

func TestShipmentService_GetMyShipments_DefaultPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipmentRepo := mock_repository.NewMockShipmentRepository(ctrl)
	shipmentService := service.NewShipmentService(mockShipmentRepo)

	ctx := context.Background()
	senderID := int64(1)
	req := &v1.ShipmentListRequest{
		Page:     0,
		PageSize: 0,
	}

	mockShipmentRepo.EXPECT().ListBySenderID(ctx, senderID, "", 1, 20).
		Return([]*model.Shipment{}, int64(0), nil)

	data, err := shipmentService.GetMyShipments(ctx, senderID, req)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 1, data.Pagination.Page)
	assert.Equal(t, 20, data.Pagination.PageSize)
}

func TestShipmentService_CalculateFreight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipmentRepo := mock_repository.NewMockShipmentRepository(ctrl)
	shipmentService := service.NewShipmentService(mockShipmentRepo)

	ctx := context.Background()
	req := &v1.CalculateFreightRequest{
		CourierCompany: "顺丰速运",
		Weight:         2.0,
		Volume:         0.05,
	}

	data, err := shipmentService.CalculateFreight(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Greater(t, data.Freight, 0.0)
	assert.Greater(t, data.CalculationDetails.BaseFreight, 0.0)
}

func TestShipmentService_CalculateFreight_SFExpress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipmentRepo := mock_repository.NewMockShipmentRepository(ctrl)
	shipmentService := service.NewShipmentService(mockShipmentRepo)

	ctx := context.Background()
	req := &v1.CalculateFreightRequest{
		CourierCompany: "顺丰速运",
		Weight:         1.0,
		Volume:         0.01,
	}

	data, err := shipmentService.CalculateFreight(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 15.0, data.CalculationDetails.BaseFreight)
}

func TestShipmentService_CalculateFreight_EMS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipmentRepo := mock_repository.NewMockShipmentRepository(ctrl)
	shipmentService := service.NewShipmentService(mockShipmentRepo)

	ctx := context.Background()
	req := &v1.CalculateFreightRequest{
		CourierCompany: "邮政EMS",
		Weight:         1.0,
		Volume:         0.01,
	}

	data, err := shipmentService.CalculateFreight(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 12.0, data.CalculationDetails.BaseFreight)
}

func TestShipmentService_GetCourierCompanies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShipmentRepo := mock_repository.NewMockShipmentRepository(ctrl)
	shipmentService := service.NewShipmentService(mockShipmentRepo)

	ctx := context.Background()

	companies := shipmentService.GetCourierCompanies(ctx)

	assert.NotNil(t, companies)
	assert.Greater(t, len(companies), 0)
	assert.Contains(t, companies, "顺丰速运")
	assert.Contains(t, companies, "中通快递")
}
