package repository

import (
	"context"

	"backend/internal/model"
)

type CourierPriceRepository interface {
	Create(ctx context.Context, price *model.CourierPrice) error
	Update(ctx context.Context, price *model.CourierPrice) error
	GetByID(ctx context.Context, id int64) (*model.CourierPrice, error)
	GetByCompany(ctx context.Context, company string) (*model.CourierPrice, error)
	List(ctx context.Context, activeOnly bool) ([]*model.CourierPrice, error)
}

type courierPriceRepository struct {
	*Repository
}

func NewCourierPriceRepository(r *Repository) CourierPriceRepository {
	return &courierPriceRepository{
		Repository: r,
	}
}

func (r *courierPriceRepository) Create(ctx context.Context, price *model.CourierPrice) error {
	return r.DB(ctx).Create(price).Error
}

func (r *courierPriceRepository) Update(ctx context.Context, price *model.CourierPrice) error {
	return r.DB(ctx).Save(price).Error
}

func (r *courierPriceRepository) GetByID(ctx context.Context, id int64) (*model.CourierPrice, error) {
	var price model.CourierPrice
	if err := r.DB(ctx).First(&price, id).Error; err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *courierPriceRepository) GetByCompany(ctx context.Context, company string) (*model.CourierPrice, error) {
	var price model.CourierPrice
	if err := r.DB(ctx).Where("courier_company = ? AND is_active = ?", company, true).First(&price).Error; err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *courierPriceRepository) List(ctx context.Context, activeOnly bool) ([]*model.CourierPrice, error) {
	var prices []*model.CourierPrice
	query := r.DB(ctx).Model(&model.CourierPrice{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	if err := query.Order("courier_company ASC").Find(&prices).Error; err != nil {
		return nil, err
	}
	return prices, nil
}
