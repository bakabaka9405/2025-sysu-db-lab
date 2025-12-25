package repository

import (
	"context"
	"time"

	"backend/internal/model"
)

type OverdueRepository interface {
	Create(ctx context.Context, record *model.OverdueRecord) error
	Update(ctx context.Context, record *model.OverdueRecord) error
	GetByID(ctx context.Context, id int64) (*model.OverdueRecord, error)
	GetByParcelID(ctx context.Context, parcelID int64) (*model.OverdueRecord, error)
	List(ctx context.Context, status string, page, pageSize int) ([]*model.OverdueRecord, int64, error)
	ListPendingWithParcels(ctx context.Context) ([]*model.OverdueRecord, error)
	UpdateReminderInfo(ctx context.Context, id int64) error
}

type overdueRepository struct {
	*Repository
}

func NewOverdueRepository(r *Repository) OverdueRepository {
	return &overdueRepository{
		Repository: r,
	}
}

func (r *overdueRepository) Create(ctx context.Context, record *model.OverdueRecord) error {
	return r.DB(ctx).Create(record).Error
}

func (r *overdueRepository) Update(ctx context.Context, record *model.OverdueRecord) error {
	return r.DB(ctx).Save(record).Error
}

func (r *overdueRepository) GetByID(ctx context.Context, id int64) (*model.OverdueRecord, error) {
	var record model.OverdueRecord
	if err := r.DB(ctx).Preload("Parcel").Preload("Parcel.Shelf").First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *overdueRepository) GetByParcelID(ctx context.Context, parcelID int64) (*model.OverdueRecord, error) {
	var record model.OverdueRecord
	if err := r.DB(ctx).Where("parcel_id = ?", parcelID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *overdueRepository) List(ctx context.Context, status string, page, pageSize int) ([]*model.OverdueRecord, int64, error) {
	var records []*model.OverdueRecord
	var total int64

	query := r.DB(ctx).Model(&model.OverdueRecord{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Parcel").Preload("Parcel.Shelf").
		Order("overdue_start_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *overdueRepository) ListPendingWithParcels(ctx context.Context) ([]*model.OverdueRecord, error) {
	var records []*model.OverdueRecord
	if err := r.DB(ctx).
		Preload("Parcel").
		Where("status = ?", model.OverdueStatusPending).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *overdueRepository) UpdateReminderInfo(ctx context.Context, id int64) error {
	now := time.Now()
	return r.DB(ctx).
		Model(&model.OverdueRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"reminder_count":   r.DB(ctx).Raw("reminder_count + 1"),
			"last_reminder_at": now,
			"status":           model.OverdueStatusReminded,
		}).Error
}
