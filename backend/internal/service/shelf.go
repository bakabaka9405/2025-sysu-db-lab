package service

import (
	"context"
	"errors"

	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/repository"
)

type ShelfService interface {
	Create(ctx context.Context, req *v1.CreateShelfRequest) (*model.Shelf, error)
	Update(ctx context.Context, id int64, req *v1.UpdateShelfRequest) (*model.Shelf, error)
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.Shelf, error)
	List(ctx context.Context, req *v1.ShelfListRequest) (*v1.ShelfListData, error)
	GetStatistics(ctx context.Context) (*v1.ShelfStatistics, error)
}

type shelfService struct {
	shelfRepo repository.ShelfRepository
}

func NewShelfService(shelfRepo repository.ShelfRepository) ShelfService {
	return &shelfService{
		shelfRepo: shelfRepo,
	}
}

func (s *shelfService) Create(ctx context.Context, req *v1.CreateShelfRequest) (*model.Shelf, error) {
	existing, _ := s.shelfRepo.GetByShelfCode(ctx, req.ShelfCode)
	if existing != nil {
		return nil, errors.New("货架编码已存在")
	}

	shelf := &model.Shelf{
		ShelfCode: req.ShelfCode,
		Area:      req.Area,
		Floor:     req.Floor,
		Column:    req.Column,
		Type:      model.ShelfType(req.Type),
		Capacity:  req.Capacity,
		Status:    model.ShelfStatusActive,
	}

	if err := s.shelfRepo.Create(ctx, shelf); err != nil {
		return nil, err
	}

	return shelf, nil
}

func (s *shelfService) Update(ctx context.Context, id int64, req *v1.UpdateShelfRequest) (*model.Shelf, error) {
	shelf, err := s.shelfRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("货架不存在")
	}

	if req.Area != "" {
		shelf.Area = req.Area
	}
	if req.Floor > 0 {
		shelf.Floor = req.Floor
	}
	if req.Column > 0 {
		shelf.Column = req.Column
	}
	if req.Type != "" {
		shelf.Type = model.ShelfType(req.Type)
	}
	if req.Capacity > 0 {
		if req.Capacity < shelf.CurrentCount {
			return nil, errors.New("新容量不能小于当前存放数量")
		}
		shelf.Capacity = req.Capacity
	}
	if req.Status != "" {
		shelf.Status = model.ShelfStatus(req.Status)
	}

	if err := s.shelfRepo.Update(ctx, shelf); err != nil {
		return nil, err
	}

	return shelf, nil
}

func (s *shelfService) Delete(ctx context.Context, id int64) error {
	shelf, err := s.shelfRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("货架不存在")
	}

	if shelf.CurrentCount > 0 {
		return errors.New("货架上还有包裹，无法删除")
	}

	return s.shelfRepo.Delete(ctx, id)
}

func (s *shelfService) GetByID(ctx context.Context, id int64) (*model.Shelf, error) {
	return s.shelfRepo.GetByID(ctx, id)
}

func (s *shelfService) List(ctx context.Context, req *v1.ShelfListRequest) (*v1.ShelfListData, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	shelves, total, err := s.shelfRepo.List(ctx, req.Area, req.Type, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &v1.ShelfListData{
		List:       convertShelvesToInfo(shelves),
		Pagination: v1.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      int(total),
			TotalPages: (int(total) + req.PageSize - 1) / req.PageSize,
		},
	}, nil
}

func (s *shelfService) GetStatistics(ctx context.Context) (*v1.ShelfStatistics, error) {
	stats := &v1.ShelfStatistics{}

	totalCapacity, totalUsed, err := s.shelfRepo.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	activeCount, _ := s.shelfRepo.CountActive(ctx)
	shelves, total, _ := s.shelfRepo.List(ctx, "", "", "", 1, 1000)

	stats.TotalShelves = int(total)
	stats.ActiveShelves = int(activeCount)
	stats.TotalCapacity = totalCapacity
	stats.TotalUsed = totalUsed
	if totalCapacity > 0 {
		stats.UtilizationRate = float64(totalUsed) / float64(totalCapacity) * 100
	}

	typeStats := make(map[string]struct {
		Count    int
		Capacity int
		Used     int
	})

	for _, shelf := range shelves {
		if shelf.Status == model.ShelfStatusActive {
			ts := typeStats[string(shelf.Type)]
			ts.Count++
			ts.Capacity += shelf.Capacity
			ts.Used += shelf.CurrentCount
			typeStats[string(shelf.Type)] = ts
		}
	}

	for shelfType, ts := range typeStats {
		var rate float64
		if ts.Capacity > 0 {
			rate = float64(ts.Used) / float64(ts.Capacity) * 100
		}
		stats.ByType = append(stats.ByType, struct {
			Type            string  `json:"type"`
			Count           int     `json:"count"`
			Capacity        int     `json:"capacity"`
			Used            int     `json:"used"`
			UtilizationRate float64 `json:"utilization_rate"`
		}{
			Type:            shelfType,
			Count:           ts.Count,
			Capacity:        ts.Capacity,
			Used:            ts.Used,
			UtilizationRate: rate,
		})
	}

	return stats, nil
}

func convertShelvesToInfo(shelves []*model.Shelf) []v1.ShelfInfo {
	infos := make([]v1.ShelfInfo, len(shelves))
	for i, s := range shelves {
		infos[i] = ConvertShelfToInfo(s)
	}
	return infos
}

func ConvertShelfToInfo(s *model.Shelf) v1.ShelfInfo {
	return v1.ShelfInfo{
		ID:           s.ID,
		ShelfCode:    s.ShelfCode,
		Area:         s.Area,
		Floor:        s.Floor,
		Column:       s.Column,
		Type:         string(s.Type),
		Capacity:     s.Capacity,
		CurrentCount: s.CurrentCount,
		Status:       string(s.Status),
	}
}
