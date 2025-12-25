package repository

import (
	"context"
	"time"

	"backend/internal/model"
)

type ParcelRepository interface {
	Create(ctx context.Context, parcel *model.Parcel) error
	Update(ctx context.Context, parcel *model.Parcel) error
	GetByID(ctx context.Context, id int64) (*model.Parcel, error)
	GetByPickupCode(ctx context.Context, pickupCode string) (*model.Parcel, error)
	GetByTrackingNumber(ctx context.Context, trackingNumber string) (*model.Parcel, error)
	List(ctx context.Context, status, courierCompany, recipientPhone string, page, pageSize int) ([]*model.Parcel, int64, error)
	ListByRecipientPhone(ctx context.Context, phone, status string, page, pageSize int) ([]*model.Parcel, int64, error)
	CountByStatus(ctx context.Context) (map[string]int, error)
	CountTodayReceived(ctx context.Context) (int, error)
	CountTodayPickedUp(ctx context.Context) (int, error)
	CountOverdue(ctx context.Context) (int, error)
	GetRecentTrends(ctx context.Context, days int) ([]map[string]interface{}, error)
	FindOverdueParcels(ctx context.Context) ([]*model.Parcel, error)
}

type parcelRepository struct {
	*Repository
}

func NewParcelRepository(r *Repository) ParcelRepository {
	return &parcelRepository{
		Repository: r,
	}
}

func (r *parcelRepository) Create(ctx context.Context, parcel *model.Parcel) error {
	return r.DB(ctx).Create(parcel).Error
}

func (r *parcelRepository) Update(ctx context.Context, parcel *model.Parcel) error {
	return r.DB(ctx).Save(parcel).Error
}

func (r *parcelRepository) GetByID(ctx context.Context, id int64) (*model.Parcel, error) {
	var parcel model.Parcel
	if err := r.DB(ctx).Preload("Shelf").First(&parcel, id).Error; err != nil {
		return nil, err
	}
	return &parcel, nil
}

func (r *parcelRepository) GetByPickupCode(ctx context.Context, pickupCode string) (*model.Parcel, error) {
	var parcel model.Parcel
	if err := r.DB(ctx).Preload("Shelf").Where("pickup_code = ?", pickupCode).First(&parcel).Error; err != nil {
		return nil, err
	}
	return &parcel, nil
}

func (r *parcelRepository) GetByTrackingNumber(ctx context.Context, trackingNumber string) (*model.Parcel, error) {
	var parcel model.Parcel
	if err := r.DB(ctx).Preload("Shelf").Where("tracking_number = ?", trackingNumber).First(&parcel).Error; err != nil {
		return nil, err
	}
	return &parcel, nil
}

func (r *parcelRepository) List(ctx context.Context, status, courierCompany, recipientPhone string, page, pageSize int) ([]*model.Parcel, int64, error) {
	var parcels []*model.Parcel
	var total int64

	db := r.DB(ctx).Model(&model.Parcel{})

	if status != "" {
		db = db.Where("status = ?", status)
	}
	if courierCompany != "" {
		db = db.Where("courier_company = ?", courierCompany)
	}
	if recipientPhone != "" {
		db = db.Where("recipient_phone LIKE ?", "%"+recipientPhone+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Shelf").Order("received_at DESC").Offset(offset).Limit(pageSize).Find(&parcels).Error; err != nil {
		return nil, 0, err
	}

	return parcels, total, nil
}

func (r *parcelRepository) ListByRecipientPhone(ctx context.Context, phone, status string, page, pageSize int) ([]*model.Parcel, int64, error) {
	var parcels []*model.Parcel
	var total int64

	db := r.DB(ctx).Model(&model.Parcel{}).Where("recipient_phone = ?", phone)

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Shelf").Order("received_at DESC").Offset(offset).Limit(pageSize).Find(&parcels).Error; err != nil {
		return nil, 0, err
	}

	return parcels, total, nil
}

func (r *parcelRepository) CountByStatus(ctx context.Context) (map[string]int, error) {
	type Result struct {
		Status string
		Count  int
	}
	var results []Result
	err := r.DB(ctx).
		Model(&model.Parcel{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

func (r *parcelRepository) CountTodayReceived(ctx context.Context) (int, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.DB(ctx).
		Model(&model.Parcel{}).
		Where("received_at >= ?", today).
		Count(&count).Error
	return int(count), err
}

func (r *parcelRepository) CountTodayPickedUp(ctx context.Context) (int, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.DB(ctx).
		Model(&model.Parcel{}).
		Where("picked_up_at >= ? AND picked_up_at IS NOT NULL", today).
		Count(&count).Error
	return int(count), err
}

func (r *parcelRepository) CountOverdue(ctx context.Context) (int, error) {
	var count int64
	err := r.DB(ctx).
		Model(&model.Parcel{}).
		Where("status = ?", model.ParcelStatusOverdue).
		Count(&count).Error
	return int(count), err
}

func (r *parcelRepository) GetRecentTrends(ctx context.Context, days int) ([]map[string]any, error) {
	type Result struct {
		Date     string
		Received int
		PickedUp int
	}
	var results []Result

	startDate := time.Now().AddDate(0, 0, -days+1).Truncate(24 * time.Hour)

	query := `
		WITH RECURSIVE dates AS (
			SELECT ? AS date
			UNION ALL
			SELECT date(date, '+1 day')
			FROM dates
			WHERE date < date('now', 'localtime')
		)
		SELECT
			dates.date,
			COALESCE(SUM(CASE WHEN DATE(p.received_at) = dates.date THEN 1 ELSE 0 END), 0) as received,
			COALESCE(SUM(CASE WHEN DATE(p.picked_up_at) = dates.date THEN 1 ELSE 0 END), 0) as picked_up
		FROM dates
		LEFT JOIN parcels p ON DATE(p.received_at) = dates.date OR DATE(p.picked_up_at) = dates.date
		GROUP BY dates.date
		ORDER BY dates.date ASC
	`

	err := r.DB(ctx).Raw(query, startDate.Format("2006-01-02")).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	trends := make([]map[string]any, len(results))
	for i, r := range results {
		trends[i] = map[string]any{
			"date":      r.Date,
			"received":  r.Received,
			"picked_up": r.PickedUp,
		}
	}

	return trends, nil
}

func (r *parcelRepository) FindOverdueParcels(ctx context.Context) ([]*model.Parcel, error) {
	var parcels []*model.Parcel
	err := r.DB(ctx).
		Where("expected_overdue_at < ? AND expected_overdue_at IS NOT NULL", time.Now()).
		Where("status NOT IN ?", []string{string(model.ParcelStatusPickedUp), string(model.ParcelStatusReturned), string(model.ParcelStatusOverdue)}).
		Find(&parcels).Error
	return parcels, err
}
