package repository

import (
	"context"

	"backend/internal/model"
	"gorm.io/gorm"
)

type ShelfRepository interface {
	Create(ctx context.Context, shelf *model.Shelf) error
	Update(ctx context.Context, shelf *model.Shelf) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.Shelf, error)
	GetByShelfCode(ctx context.Context, shelfCode string) (*model.Shelf, error)
	FindAvailableShelfByType(ctx context.Context, shelfType model.ShelfType) (*model.Shelf, error)
	IncrementCount(ctx context.Context, id int64) error
	DecrementCount(ctx context.Context, id int64) error
	GetStatistics(ctx context.Context) (totalCapacity, totalUsed int, err error)
	List(ctx context.Context, area, shelfType, status string, page, pageSize int) ([]*model.Shelf, int64, error)
	CountByType(ctx context.Context) ([]map[string]interface{}, error)
	CountActive(ctx context.Context) (int64, error)
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

func (r *shelfRepository) Delete(ctx context.Context, id int64) error {
	return r.DB(ctx).Delete(&model.Shelf{}, id).Error
}

func (r *shelfRepository) List(ctx context.Context, area, shelfType, status string, page, pageSize int) ([]*model.Shelf, int64, error) {
	var shelves []*model.Shelf
	var total int64

	query := r.DB(ctx).Model(&model.Shelf{})
	if area != "" {
		query = query.Where("area = ?", area)
	}
	if shelfType != "" {
		query = query.Where("type = ?", shelfType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("shelf_code ASC").Offset(offset).Limit(pageSize).Find(&shelves).Error; err != nil {
		return nil, 0, err
	}

	return shelves, total, nil
}

func (r *shelfRepository) CountByType(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.DB(ctx).
		Model(&model.Shelf{}).
		Select("type, COUNT(*) as count, SUM(capacity) as capacity, SUM(current_count) as used").
		Where("status = ?", model.ShelfStatusActive).
		Group("type").
		Find(&results).Error
	return results, err
}

func (r *shelfRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB(ctx).Model(&model.Shelf{}).Where("status = ?", model.ShelfStatusActive).Count(&count).Error
	return count, err
}
