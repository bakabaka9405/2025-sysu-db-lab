package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/service"
	mock_repository "backend/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestParcelService_ReceiveParcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()
	req := &v1.ReceiveParcelRequest{
		TrackingNumber:  "TEST123456",
		RecipientName:   "Test User",
		RecipientPhone:  "13800138000",
		RecipientIDCard: "440101199001011234",
		CourierCompany:  "SF Express",
		Size:            "small",
		Weight:          1.5,
		Notes:           "Test notes",
	}

	mockParcelRepo.EXPECT().GetByTrackingNumber(ctx, req.TrackingNumber).Return(nil, nil)
	mockParcelRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	mockParcelRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&model.Parcel{
		ID:             1,
		TrackingNumber: req.TrackingNumber,
		PickupCode:     "A12345",
		RecipientName:  req.RecipientName,
		RecipientPhone: req.RecipientPhone,
		Status:         model.ParcelStatusReceived,
	}, nil)

	parcel, err := parcelService.ReceiveParcel(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, parcel)
	assert.Equal(t, req.TrackingNumber, parcel.TrackingNumber)
	assert.Equal(t, model.ParcelStatusReceived, parcel.Status)
}

func TestParcelService_ReceiveParcel_DuplicateTracking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()
	req := &v1.ReceiveParcelRequest{
		TrackingNumber: "TEST123456",
		RecipientName:  "Test User",
		RecipientPhone: "13800138000",
		CourierCompany: "SF Express",
		Size:           "small",
		Weight:         1.5,
	}

	mockParcelRepo.EXPECT().GetByTrackingNumber(ctx, req.TrackingNumber).Return(&model.Parcel{
		ID: 1,
	}, nil)

	_, err := parcelService.ReceiveParcel(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "快递单号已存在")
}

func TestParcelService_PickupParcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()
	req := &v1.PickupParcelRequest{
		PickupCode:     "A12345",
		RecipientPhone: "13800138000",
	}

	shelfID := int64(1)
	mockParcel := &model.Parcel{
		ID:             1,
		PickupCode:     req.PickupCode,
		RecipientPhone: req.RecipientPhone,
		Status:         model.ParcelStatusReady,
		ShelfID:        &shelfID,
	}

	mockParcelRepo.EXPECT().GetByPickupCode(ctx, req.PickupCode).Return(mockParcel, nil)
	mockTm.EXPECT().Transaction(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			mockParcelRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
			mockShelfRepo.EXPECT().DecrementCount(ctx, shelfID).Return(nil)
			return fn(ctx)
		},
	)

	parcel, err := parcelService.PickupParcel(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, parcel)
	assert.Equal(t, model.ParcelStatusPickedUp, parcel.Status)
}

func TestParcelService_PickupParcel_InvalidPickupCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()
	req := &v1.PickupParcelRequest{
		PickupCode:     "INVALID",
		RecipientPhone: "13800138000",
	}

	mockParcelRepo.EXPECT().GetByPickupCode(ctx, req.PickupCode).Return(nil, errors.New("not found"))

	_, err := parcelService.PickupParcel(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "取件码不存在")
}

func TestParcelService_PickupParcel_PhoneMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()
	req := &v1.PickupParcelRequest{
		PickupCode:     "A12345",
		RecipientPhone: "13800138000",
	}

	mockParcel := &model.Parcel{
		ID:             1,
		PickupCode:     req.PickupCode,
		RecipientPhone: "13900139000",
		Status:         model.ParcelStatusReceived,
	}

	mockParcelRepo.EXPECT().GetByPickupCode(ctx, req.PickupCode).Return(mockParcel, nil)

	_, err := parcelService.PickupParcel(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "手机号不匹配")
}

func TestParcelService_PickupParcel_AlreadyPickedUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()
	req := &v1.PickupParcelRequest{
		PickupCode:     "A12345",
		RecipientPhone: "13800138000",
	}

	mockParcel := &model.Parcel{
		ID:             1,
		PickupCode:     req.PickupCode,
		RecipientPhone: req.RecipientPhone,
		Status:         model.ParcelStatusPickedUp,
	}

	mockParcelRepo.EXPECT().GetByPickupCode(ctx, req.PickupCode).Return(mockParcel, nil)

	_, err := parcelService.PickupParcel(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已被取走")
}

func TestParcelService_GetParcelByPickupCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()
	pickupCode := "A12345"

	mockParcel := &model.Parcel{
		ID:         1,
		PickupCode: pickupCode,
		Status:     model.ParcelStatusReceived,
	}

	mockParcelRepo.EXPECT().GetByPickupCode(ctx, pickupCode).Return(mockParcel, nil)

	parcel, err := parcelService.GetParcelByPickupCode(ctx, pickupCode)

	assert.NoError(t, err)
	assert.NotNil(t, parcel)
	assert.Equal(t, pickupCode, parcel.PickupCode)
}

func TestParcelService_GetParcelByTrackingNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()
	trackingNumber := "TEST123456"

	mockParcel := &model.Parcel{
		ID:             1,
		TrackingNumber: trackingNumber,
		Status:         model.ParcelStatusReceived,
	}

	mockParcelRepo.EXPECT().GetByTrackingNumber(ctx, trackingNumber).Return(mockParcel, nil)

	parcel, err := parcelService.GetParcelByTrackingNumber(ctx, trackingNumber)

	assert.NoError(t, err)
	assert.NotNil(t, parcel)
	assert.Equal(t, trackingNumber, parcel.TrackingNumber)
}

func TestParcelService_GetParcelList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()
	req := &v1.ParcelListRequest{
		Page:     1,
		PageSize: 20,
		Status:   "received",
	}

	mockParcels := []*model.Parcel{
		{
			ID:         1,
			PickupCode: "A12345",
			Status:     model.ParcelStatusReceived,
			ReceivedAt: time.Now(),
		},
	}

	mockParcelRepo.EXPECT().List(ctx, req.Status, req.CourierCompany, req.RecipientPhone, req.Page, req.PageSize).
		Return(mockParcels, int64(1), nil)

	data, err := parcelService.GetParcelList(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 1, len(data.List))
	assert.Equal(t, 1, data.Pagination.Total)
}

func TestParcelService_GetDashboardStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParcelRepo := mock_repository.NewMockParcelRepository(ctrl)
	mockShelfRepo := mock_repository.NewMockShelfRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)

	parcelService := service.NewParcelService(mockParcelRepo, mockShelfRepo, mockTm)

	ctx := context.Background()

	mockParcelRepo.EXPECT().CountTodayReceived(ctx).Return(10, nil)
	mockParcelRepo.EXPECT().CountTodayPickedUp(ctx).Return(5, nil)
	mockParcelRepo.EXPECT().CountOverdue(ctx).Return(2, nil)
	mockParcelRepo.EXPECT().CountByStatus(ctx).Return(map[string]int{
		"received":  15,
		"picked_up": 10,
	}, nil)
	mockShelfRepo.EXPECT().GetStatistics(ctx).Return(1000, 500, nil)
	mockParcelRepo.EXPECT().GetRecentTrends(ctx, 7).Return([]map[string]interface{}{
		{"date": "2024-01-01", "received": int64(10), "picked_up": int64(5)},
	}, nil)

	stats, err := parcelService.GetDashboardStatistics(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 10, stats.Today.ReceivedCount)
	assert.Equal(t, 5, stats.Today.PickedUpCount)
	assert.Equal(t, 2, stats.Today.OverdueCount)
	assert.Equal(t, 1000, stats.Shelves.TotalCapacity)
	assert.Equal(t, 500, stats.Shelves.TotalUsed)
	assert.Equal(t, 50.0, stats.Shelves.UtilizationRate)
}
