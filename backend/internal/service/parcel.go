package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/repository"
)

type ParcelService interface {
	ReceiveParcel(ctx context.Context, req *v1.ReceiveParcelRequest) (*model.Parcel, error)
	ShelveParcel(ctx context.Context, id int64) (*model.Parcel, error)
	PickupParcel(ctx context.Context, req *v1.PickupParcelRequest) (*model.Parcel, error)
	GetParcelByPickupCode(ctx context.Context, pickupCode string) (*model.Parcel, error)
	GetParcelByTrackingNumber(ctx context.Context, trackingNumber string) (*model.Parcel, error)
	GetParcelList(ctx context.Context, req *v1.ParcelListRequest) (*v1.ParcelListData, error)
	GetMyParcels(ctx context.Context, phone string, req *v1.ParcelListRequest) (*v1.ParcelListData, error)
	GetDashboardStatistics(ctx context.Context) (*v1.DashboardStatistics, error)
}

type parcelService struct {
	parcelRepo repository.ParcelRepository
	shelfRepo  repository.ShelfRepository
	tm         repository.Transaction
}

func NewParcelService(
	parcelRepo repository.ParcelRepository,
	shelfRepo repository.ShelfRepository,
	tm repository.Transaction,
) ParcelService {
	return &parcelService{
		parcelRepo: parcelRepo,
		shelfRepo:  shelfRepo,
		tm:         tm,
	}
}

func (s *parcelService) ReceiveParcel(ctx context.Context, req *v1.ReceiveParcelRequest) (*model.Parcel, error) {
	existing, _ := s.parcelRepo.GetByTrackingNumber(ctx, req.TrackingNumber)
	if existing != nil {
		return nil, errors.New("快递单号已存在")
	}

	// 入库时只记录包裹信息，状态为 received，不分配货架
	parcel := &model.Parcel{
		TrackingNumber:  req.TrackingNumber,
		PickupCode:      generatePickupCode(),
		RecipientName:   req.RecipientName,
		RecipientPhone:  req.RecipientPhone,
		RecipientIDCard: req.RecipientIDCard,
		CourierCompany:  req.CourierCompany,
		Size:            model.ParcelSize(req.Size),
		Weight:          req.Weight,
		Status:          model.ParcelStatusReceived,
		ReceivedAt:      time.Now(),
		Notes:           req.Notes,
	}

	if err := s.parcelRepo.Create(ctx, parcel); err != nil {
		return nil, err
	}

	return s.parcelRepo.GetByID(ctx, parcel.ID)
}

// ShelveParcel 包裹上架
func (s *parcelService) ShelveParcel(ctx context.Context, id int64) (*model.Parcel, error) {
	parcel, err := s.parcelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("包裹不存在")
	}

	if parcel.Status != model.ParcelStatusReceived {
		return nil, errors.New("只有已入库状态的包裹才能上架")
	}

	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// 分配货架
		shelfType := model.ShelfType(parcel.Size)
		shelf, err := s.shelfRepo.FindAvailableShelfByType(ctx, shelfType)
		if err != nil {
			return errors.New("暂无可用货架")
		}

		// 更新包裹状态
		now := time.Now()
		parcel.Status = model.ParcelStatusReadyForPickup
		parcel.ShelfID = &shelf.ID
		parcel.ShelvedAt = &now

		// 设置预计滞留时间（上架后3天）
		overdueTime := now.Add(3 * 24 * time.Hour)
		parcel.ExpectedOverdueAt = &overdueTime

		if err := s.parcelRepo.Update(ctx, parcel); err != nil {
			return err
		}

		if err := s.shelfRepo.IncrementCount(ctx, shelf.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.parcelRepo.GetByID(ctx, parcel.ID)
}

func (s *parcelService) PickupParcel(ctx context.Context, req *v1.PickupParcelRequest) (*model.Parcel, error) {
	parcel, err := s.parcelRepo.GetByPickupCode(ctx, req.PickupCode)
	if err != nil {
		return nil, errors.New("取件码不存在")
	}

	if parcel.RecipientPhone != req.RecipientPhone {
		return nil, errors.New("手机号不匹配")
	}

	if parcel.Status == model.ParcelStatusPickedUp {
		return nil, errors.New("该包裹已被取走")
	}

	if parcel.Status == model.ParcelStatusReturned {
		return nil, errors.New("该包裹已被退回")
	}

	// 只有待取件或滞留状态的包裹才能取件
	if parcel.Status != model.ParcelStatusReadyForPickup && parcel.Status != model.ParcelStatusOverdue {
		return nil, errors.New("该包裹当前状态不允许取件")
	}

	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		now := time.Now()
		parcel.Status = model.ParcelStatusPickedUp
		parcel.PickedUpAt = &now

		if err := s.parcelRepo.Update(ctx, parcel); err != nil {
			return err
		}

		if parcel.ShelfID != nil {
			if err := s.shelfRepo.DecrementCount(ctx, *parcel.ShelfID); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return parcel, nil
}

func (s *parcelService) GetParcelByPickupCode(ctx context.Context, pickupCode string) (*model.Parcel, error) {
	return s.parcelRepo.GetByPickupCode(ctx, pickupCode)
}

func (s *parcelService) GetParcelByTrackingNumber(ctx context.Context, trackingNumber string) (*model.Parcel, error) {
	return s.parcelRepo.GetByTrackingNumber(ctx, trackingNumber)
}

func (s *parcelService) GetParcelList(ctx context.Context, req *v1.ParcelListRequest) (*v1.ParcelListData, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	parcels, total, err := s.parcelRepo.List(ctx, req.Status, req.CourierCompany, req.RecipientPhone, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &v1.ParcelListData{
		List: convertParcelsToInfo(parcels),
		Pagination: v1.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      int(total),
			TotalPages: (int(total) + req.PageSize - 1) / req.PageSize,
		},
	}, nil
}

func (s *parcelService) GetMyParcels(ctx context.Context, phone string, req *v1.ParcelListRequest) (*v1.ParcelListData, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	parcels, total, err := s.parcelRepo.ListByRecipientPhone(ctx, phone, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &v1.ParcelListData{
		List: convertParcelsToInfo(parcels),
		Pagination: v1.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      int(total),
			TotalPages: (int(total) + req.PageSize - 1) / req.PageSize,
		},
	}, nil
}

func (s *parcelService) GetDashboardStatistics(ctx context.Context) (*v1.DashboardStatistics, error) {
	stats := &v1.DashboardStatistics{}

	todayReceived, _ := s.parcelRepo.CountTodayReceived(ctx)
	todayPickedUp, _ := s.parcelRepo.CountTodayPickedUp(ctx)
	overdueCount, _ := s.parcelRepo.CountOverdue(ctx)

	stats.Today.ReceivedCount = todayReceived
	stats.Today.PickedUpCount = todayPickedUp
	stats.Today.OverdueCount = overdueCount

	statusCounts, _ := s.parcelRepo.CountByStatus(ctx)
	stats.Parcels.Received = statusCounts[string(model.ParcelStatusReceived)]
	stats.Parcels.ReadyForPickup = statusCounts[string(model.ParcelStatusReadyForPickup)]
	stats.Parcels.PickedUp = statusCounts[string(model.ParcelStatusPickedUp)]
	stats.Parcels.Overdue = statusCounts[string(model.ParcelStatusOverdue)]
	stats.Parcels.Returned = statusCounts[string(model.ParcelStatusReturned)]
	stats.Parcels.Total = 0
	for _, count := range statusCounts {
		stats.Parcels.Total += count
	}

	totalCapacity, totalUsed, _ := s.shelfRepo.GetStatistics(ctx)
	stats.Shelves.TotalCapacity = totalCapacity
	stats.Shelves.TotalUsed = totalUsed
	if totalCapacity > 0 {
		stats.Shelves.UtilizationRate = float64(totalUsed) / float64(totalCapacity) * 100
	}

	trends, _ := s.parcelRepo.GetRecentTrends(ctx, 7)
	stats.RecentTrends = make([]struct {
		Date     string `json:"date"`
		Received int    `json:"received"`
		PickedUp int    `json:"picked_up"`
	}, len(trends))
	for i, trend := range trends {
		stats.RecentTrends[i].Date = trend["date"].(string)
		stats.RecentTrends[i].Received = int(trend["received"].(int64))
		stats.RecentTrends[i].PickedUp = int(trend["picked_up"].(int64))
	}

	return stats, nil
}

func generatePickupCode() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := string(letters[rand.Intn(len(letters))])
	code += fmt.Sprintf("%05d", rand.Intn(90000)+10000)
	return code
}

// ConvertParcelToInfo converts a single parcel to ParcelInfo
func ConvertParcelToInfo(p *model.Parcel) v1.ParcelInfo {
	info := v1.ParcelInfo{
		ID:                p.ID,
		TrackingNumber:    p.TrackingNumber,
		PickupCode:        p.PickupCode,
		RecipientName:     p.RecipientName,
		RecipientPhone:    p.RecipientPhone,
		CourierCompany:    p.CourierCompany,
		Size:              string(p.Size),
		Weight:            p.Weight,
		Status:            string(p.Status),
		ReceivedAt:        p.ReceivedAt,
		ShelvedAt:         p.ShelvedAt,
		PickedUpAt:        p.PickedUpAt,
		ExpectedOverdueAt: p.ExpectedOverdueAt,
		Notes:             p.Notes,
	}
	if p.Shelf != nil {
		info.Shelf = &v1.ShelfInfo{
			ID:           p.Shelf.ID,
			ShelfCode:    p.Shelf.ShelfCode,
			Area:         p.Shelf.Area,
			Floor:        p.Shelf.Floor,
			Column:       p.Shelf.Column,
			Type:         string(p.Shelf.Type),
			Capacity:     p.Shelf.Capacity,
			CurrentCount: p.Shelf.CurrentCount,
			Status:       string(p.Shelf.Status),
		}
	}
	return info
}

func convertParcelsToInfo(parcels []*model.Parcel) []v1.ParcelInfo {
	infos := make([]v1.ParcelInfo, len(parcels))
	for i, p := range parcels {
		infos[i] = ConvertParcelToInfo(p)
	}
	return infos
}
