package repository

import (
	"context"

	"backend/internal/model"
	"gorm.io/gorm"
)

type ShelfRepository interface {
	Create(ctx context.Context, shelf *model.Shelf) error
	Update(ctx context.Context, shelf *model.Shelf) error
	GetByID(ctx context.Context, id int64) (*model.Shelf, error)
	GetByShelfCode(ctx context.Context, shelfCode string) (*model.Shelf, error)
	FindAvailableShelfByType(ctx context.Context, shelfType model.ShelfType) (*model.Shelf, error)
	IncrementCount(ctx context.Context, id int64) error
	DecrementCount(ctx context.Context, id int64) error
	GetStatistics(ctx context.Context) (totalCapacity, totalUsed int, err error)
}

type shelfRepository struct {
	*Repository
}

func NewShelfRepository(r *Repository) ShelfRepository {
	return &shelfRepository{
		Repository: r,
	}
}

func (r *shelfRepository) Create(ctx context.Context, shelf *model.Shelf) error {
	return r.DB(ctx).Create(shelf).Error
}

func (r *shelfRepository) Update(ctx context.Context, shelf *model.Shelf) error {
	return r.DB(ctx).Save(shelf).Error
}

func (r *shelfRepository) GetByID(ctx context.Context, id int64) (*model.Shelf, error) {
	var shelf model.Shelf
	if err := r.DB(ctx).First(&shelf, id).Error; err != nil {
		return nil, err
	}
	return &shelf, nil
}

func (r *shelfRepository) GetByShelfCode(ctx context.Context, shelfCode string) (*model.Shelf, error) {
	var shelf model.Shelf
	if err := r.DB(ctx).Where("shelf_code = ?", shelfCode).First(&shelf).Error; err != nil {
		return nil, err
	}
	return &shelf, nil
}

func (r *shelfRepository) FindAvailableShelfByType(ctx context.Context, shelfType model.ShelfType) (*model.Shelf, error) {
	var shelf model.Shelf
	err := r.DB(ctx).
		Where("type = ? AND status = ? AND current_count < capacity", shelfType, model.ShelfStatusActive).
		Order("current_count ASC").
		First(&shelf).Error
	if err != nil {
		return nil, err
	}
	return &shelf, nil
}

func (r *shelfRepository) IncrementCount(ctx context.Context, id int64) error {
	return r.DB(ctx).
		Model(&model.Shelf{}).
		Where("id = ?", id).
		UpdateColumn("current_count", gorm.Expr("current_count + ?", 1)).
		Error
}

func (r *shelfRepository) DecrementCount(ctx context.Context, id int64) error {
	return r.DB(ctx).
		Model(&model.Shelf{}).
		Where("id = ? AND current_count > ?", id, 0).
		UpdateColumn("current_count", gorm.Expr("current_count - ?", 1)).
		Error
}

func (r *shelfRepository) GetStatistics(ctx context.Context) (totalCapacity, totalUsed int, err error) {
	type Result struct {
		TotalCapacity int
		TotalUsed     int
	}
	var result Result
	err = r.DB(ctx).
		Model(&model.Shelf{}).
		Select("SUM(capacity) as total_capacity, SUM(current_count) as total_used").
		Where("status = ?", model.ShelfStatusActive).
		Scan(&result).Error
	return result.TotalCapacity, result.TotalUsed, err
}
