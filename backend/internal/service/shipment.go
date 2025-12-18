package service

import (
	"context"
	"fmt"
	"math"
	"time"

	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/repository"
)

type ShipmentService interface {
	CreateShipment(ctx context.Context, senderID int64, senderName, senderPhone string, req *v1.CreateShipmentRequest) (*model.Shipment, error)
	GetShipmentByID(ctx context.Context, id int64) (*model.Shipment, error)
	GetMyShipments(ctx context.Context, senderID int64, req *v1.ShipmentListRequest) (*v1.ShipmentListData, error)
	CalculateFreight(ctx context.Context, req *v1.CalculateFreightRequest) (*v1.CalculateFreightData, error)
	GetCourierCompanies(ctx context.Context) []string
}

type shipmentService struct {
	shipmentRepo repository.ShipmentRepository
}

func NewShipmentService(shipmentRepo repository.ShipmentRepository) ShipmentService {
	return &shipmentService{
		shipmentRepo: shipmentRepo,
	}
}

var courierCompanies = []string{
	"顺丰速运",
	"中通快递",
	"圆通速递",
	"申通快递",
	"韵达快递",
	"邮政EMS",
	"百世快递",
	"德邦物流",
}

func (s *shipmentService) CreateShipment(ctx context.Context, senderID int64, senderName, senderPhone string, req *v1.CreateShipmentRequest) (*model.Shipment, error) {
	freightData := calculateFreightInternal(req.CourierCompany, req.Weight, req.Volume)

	shipment := &model.Shipment{
		ShipmentNumber:   generateShipmentNumber(),
		SenderID:         senderID,
		SenderName:       senderName,
		SenderPhone:      senderPhone,
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		RecipientAddress: req.RecipientAddress,
		CourierCompany:   req.CourierCompany,
		Weight:           req.Weight,
		Volume:           req.Volume,
		Freight:          freightData.Freight,
		Status:           model.ShipmentStatusPending,
		CreatedAt:        time.Now(),
	}

	if err := s.shipmentRepo.Create(ctx, shipment); err != nil {
		return nil, err
	}

	return shipment, nil
}

func (s *shipmentService) GetShipmentByID(ctx context.Context, id int64) (*model.Shipment, error) {
	return s.shipmentRepo.GetByID(ctx, id)
}

func (s *shipmentService) GetMyShipments(ctx context.Context, senderID int64, req *v1.ShipmentListRequest) (*v1.ShipmentListData, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	shipments, total, err := s.shipmentRepo.ListBySenderID(ctx, senderID, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &v1.ShipmentListData{
		List:       convertShipmentsToInfo(shipments),
		Pagination: v1.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      int(total),
			TotalPages: (int(total) + req.PageSize - 1) / req.PageSize,
		},
	}, nil
}

func (s *shipmentService) CalculateFreight(ctx context.Context, req *v1.CalculateFreightRequest) (*v1.CalculateFreightData, error) {
	return calculateFreightInternal(req.CourierCompany, req.Weight, req.Volume), nil
}

func (s *shipmentService) GetCourierCompanies(ctx context.Context) []string {
	return courierCompanies
}

func calculateFreightInternal(courierCompany string, weight, volume float64) *v1.CalculateFreightData {
	var basePrice float64 = 10

	if courierCompany == "顺丰速运" {
		basePrice = 15
	} else if courierCompany == "邮政EMS" {
		basePrice = 12
	}

	additionalWeightPrice := 5.0
	additionalWeight := math.Max(0, weight-1)
	weightFreight := additionalWeight * additionalWeightPrice

	volumePrice := 100.0
	volumeFee := math.Max(0, (volume-0.01)*volumePrice)

	total := basePrice + weightFreight + volumeFee

	data := &v1.CalculateFreightData{
		Freight: math.Round(total*100) / 100,
	}
	data.CalculationDetails.BaseFreight = basePrice
	data.CalculationDetails.WeightFreight = math.Round(weightFreight*100) / 100
	data.CalculationDetails.VolumeFreight = math.Round(volumeFee*100) / 100

	return data
}

func generateShipmentNumber() string {
	year := time.Now().Year()
	timestamp := time.Now().UnixNano() / 1000000
	return fmt.Sprintf("SH%d%d", year, timestamp%100000000)
}

// ConvertShipmentToInfo converts a single shipment to ShipmentInfo
func ConvertShipmentToInfo(s *model.Shipment) v1.ShipmentInfo {
	return v1.ShipmentInfo{
		ID:               s.ID,
		ShipmentNumber:   s.ShipmentNumber,
		SenderName:       s.SenderName,
		SenderPhone:      s.SenderPhone,
		RecipientName:    s.RecipientName,
		RecipientPhone:   s.RecipientPhone,
		RecipientAddress: s.RecipientAddress,
		CourierCompany:   s.CourierCompany,
		Weight:           s.Weight,
		Volume:           s.Volume,
		Freight:          s.Freight,
		Status:           string(s.Status),
		CreatedAt:        s.CreatedAt,
		ShippedAt:        s.ShippedAt,
	}
}

func convertShipmentsToInfo(shipments []*model.Shipment) []v1.ShipmentInfo {
	infos := make([]v1.ShipmentInfo, len(shipments))
	for i, s := range shipments {
		infos[i] = ConvertShipmentToInfo(s)
	}
	return infos
}
