package repository

import (
	"context"

	"backend/internal/model"
)

type ShipmentRepository interface {
	Create(ctx context.Context, shipment *model.Shipment) error
	Update(ctx context.Context, shipment *model.Shipment) error
	GetByID(ctx context.Context, id int64) (*model.Shipment, error)
	GetByShipmentNumber(ctx context.Context, shipmentNumber string) (*model.Shipment, error)
	ListBySenderID(ctx context.Context, senderID int64, status string, page, pageSize int) ([]*model.Shipment, int64, error)
}

type shipmentRepository struct {
	*Repository
}

func NewShipmentRepository(r *Repository) ShipmentRepository {
	return &shipmentRepository{
		Repository: r,
	}
}

func (r *shipmentRepository) Create(ctx context.Context, shipment *model.Shipment) error {
	return r.DB(ctx).Create(shipment).Error
}

func (r *shipmentRepository) Update(ctx context.Context, shipment *model.Shipment) error {
	return r.DB(ctx).Save(shipment).Error
}

func (r *shipmentRepository) GetByID(ctx context.Context, id int64) (*model.Shipment, error) {
	var shipment model.Shipment
	if err := r.DB(ctx).First(&shipment, id).Error; err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) GetByShipmentNumber(ctx context.Context, shipmentNumber string) (*model.Shipment, error) {
	var shipment model.Shipment
	if err := r.DB(ctx).Where("shipment_number = ?", shipmentNumber).First(&shipment).Error; err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) ListBySenderID(ctx context.Context, senderID int64, status string, page, pageSize int) ([]*model.Shipment, int64, error) {
	var shipments []*model.Shipment
	var total int64

	db := r.DB(ctx).Model(&model.Shipment{}).Where("sender_id = ?", senderID)

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}
