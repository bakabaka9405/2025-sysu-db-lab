package service

import (
	"context"
	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponseData, error)
	GetProfile(ctx context.Context, userId int64) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId int64, req *v1.UpdateProfileRequest) error
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

type userService struct {
	userRepo repository.UserRepository
	*Service
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// 检查用户名是否已存在
	existUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if existUser != nil {
		return v1.ErrUsernameExists
	}

	// 检查手机号是否已存在
	existUser, err = s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if existUser != nil {
		return v1.ErrPhoneExists
	}

	// 如果提供了邮箱，检查是否已存在
	if req.Email != "" {
		existUser, err = s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return v1.ErrInternalServerError
		}
		if existUser != nil {
			return v1.ErrEmailExists
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户对象
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		RealName: req.RealName,
		Phone:    req.Phone,
		IDCard:   req.IDCard,
		Email:    req.Email,
		Role:     "user", // 默认角色为普通用户
	}

	// 使用事务创建用户
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponseData, error) {
	var user *model.User
	var err error

	// 判断是手机号还是用户名（简单判断：纯数字认为是手机号）
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if phoneRegex.MatchString(req.UsernameOrPhone) {
		// 使用手机号登录
		user, err = s.userRepo.GetByPhone(ctx, req.UsernameOrPhone)
	} else {
		// 使用用户名登录
		user, err = s.userRepo.GetByUsername(ctx, req.UsernameOrPhone)
	}

	if err != nil || user == nil {
		return nil, v1.ErrUnauthorized
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, v1.ErrUnauthorized
	}

	// 生成JWT token (使用用户ID和角色)
	token, err := s.jwt.GenToken(fmt.Sprintf("%d", user.ID), user.Role, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return nil, err
	}

	return &v1.LoginResponseData{
		AccessToken: token,
		User: v1.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			RealName: user.RealName,
			Phone:    user.Phone,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, userId int64) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		ID:       user.ID,
		Username: user.Username,
		RealName: user.RealName,
		Phone:    user.Phone,
		IDCard:   user.IDCard,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId int64, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.RealName = req.RealName
	user.Email = req.Email

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
