package service

import (
	"context"
	"errors"
	"time"

	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/repository"
)

type OverdueService interface {
	List(ctx context.Context, req *v1.OverdueListRequest) (*v1.OverdueListData, error)
	SendReminder(ctx context.Context, id int64) error
	BatchReturn(ctx context.Context, ids []int64) (*v1.BatchResult, error)
	CheckAndMarkOverdue(ctx context.Context) (int, error)
}

type overdueService struct {
	overdueRepo repository.OverdueRepository
	parcelRepo  repository.ParcelRepository
	shelfRepo   repository.ShelfRepository
	tm          repository.Transaction
}

func NewOverdueService(
	overdueRepo repository.OverdueRepository,
	parcelRepo repository.ParcelRepository,
	shelfRepo repository.ShelfRepository,
	tm repository.Transaction,
) OverdueService {
	return &overdueService{
		overdueRepo: overdueRepo,
		parcelRepo:  parcelRepo,
		shelfRepo:   shelfRepo,
		tm:          tm,
	}
}

func (s *overdueService) List(ctx context.Context, req *v1.OverdueListRequest) (*v1.OverdueListData, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	records, total, err := s.overdueRepo.List(ctx, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &v1.OverdueListData{
		List:       s.convertToOverdueInfoList(records),
		Pagination: v1.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      int(total),
			TotalPages: (int(total) + req.PageSize - 1) / req.PageSize,
		},
	}, nil
}

func (s *overdueService) SendReminder(ctx context.Context, id int64) error {
	record, err := s.overdueRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("滞留记录不存在")
	}

	if record.Status == model.OverdueStatusReturned {
		return errors.New("该包裹已退回")
	}

	// TODO: 实际发送提醒（短信/推送）
	// 这里只更新提醒记录
	now := time.Now()
	record.ReminderCount++
	record.LastReminderAt = &now
	record.Status = model.OverdueStatusReminded

	return s.overdueRepo.Update(ctx, record)
}

func (s *overdueService) BatchReturn(ctx context.Context, ids []int64) (*v1.BatchResult, error) {
	result := &v1.BatchResult{}

	for _, id := range ids {
		err := s.tm.Transaction(ctx, func(ctx context.Context) error {
			record, err := s.overdueRepo.GetByID(ctx, id)
			if err != nil {
				return errors.New("滞留记录不存在")
			}

			if record.Status == model.OverdueStatusReturned {
				return errors.New("该包裹已退回")
			}

			// 更新包裹状态
			parcel, err := s.parcelRepo.GetByID(ctx, record.ParcelID)
			if err != nil {
				return err
			}

			parcel.Status = model.ParcelStatusReturned
			if err := s.parcelRepo.Update(ctx, parcel); err != nil {
				return err
			}

			// 更新货架计数
			if parcel.ShelfID != nil {
				if err := s.shelfRepo.DecrementCount(ctx, *parcel.ShelfID); err != nil {
					return err
				}
			}

			// 更新滞留记录状态
			record.Status = model.OverdueStatusReturned
			record.OverdueFee = record.CalculateFee()
			return s.overdueRepo.Update(ctx, record)
		})

		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Success++
		}
	}

	return result, nil
}

func (s *overdueService) CheckAndMarkOverdue(ctx context.Context) (int, error) {
	// 查找所有超期未取的包裹
	parcels, err := s.parcelRepo.FindOverdueParcels(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, parcel := range parcels {
		err := s.tm.Transaction(ctx, func(ctx context.Context) error {
			// 检查是否已有滞留记录
			existing, _ := s.overdueRepo.GetByParcelID(ctx, parcel.ID)
			if existing != nil {
				return nil // 跳过已有记录的
			}

			// 更新包裹状态
			parcel.Status = model.ParcelStatusOverdue
			if err := s.parcelRepo.Update(ctx, parcel); err != nil {
				return err
			}

			// 创建滞留记录
			record := &model.OverdueRecord{
				ParcelID:       parcel.ID,
				OverdueStartAt: *parcel.ExpectedOverdueAt,
				Status:         model.OverdueStatusPending,
			}
			return s.overdueRepo.Create(ctx, record)
		})

		if err == nil {
			count++
		}
	}

	return count, nil
}

func (s *overdueService) convertToOverdueInfoList(records []*model.OverdueRecord) []v1.OverdueInfo {
	infos := make([]v1.OverdueInfo, len(records))
	for i, r := range records {
		infos[i] = s.convertToOverdueInfo(r)
	}
	return infos
}

func (s *overdueService) convertToOverdueInfo(r *model.OverdueRecord) v1.OverdueInfo {
	info := v1.OverdueInfo{
		ID:             r.ID,
		ParcelID:       r.ParcelID,
		OverdueStartAt: r.OverdueStartAt,
		OverdueDays:    int(time.Since(r.OverdueStartAt).Hours() / 24),
		ReminderCount:  r.ReminderCount,
		LastReminderAt: r.LastReminderAt,
		OverdueFee:     r.CalculateFee(),
		Status:         string(r.Status),
	}

	if r.Parcel != nil {
		info.TrackingNumber = r.Parcel.TrackingNumber
		info.PickupCode = r.Parcel.PickupCode
		info.RecipientName = r.Parcel.RecipientName
		info.RecipientPhone = r.Parcel.RecipientPhone
		info.CourierCompany = r.Parcel.CourierCompany
		info.Size = string(r.Parcel.Size)
		if r.Parcel.Shelf != nil {
			info.ShelfCode = r.Parcel.Shelf.ShelfCode
		}
	}

	return info
}
