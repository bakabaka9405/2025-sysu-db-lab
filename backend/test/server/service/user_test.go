package service_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	v1 "backend/api/v1"
	"backend/pkg/jwt"
	"backend/test/mocks/repository"
	"os"
	"testing"

	"backend/internal/model"
	"backend/internal/service"
	"backend/pkg/config"
	"backend/pkg/log"
	"backend/pkg/sid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	logger *log.Logger
	j      *jwt.JWT
	sf     *sid.Sid
)

func TestMain(m *testing.M) {
	fmt.Println("begin")

	err := os.Setenv("APP_CONF", "../../../config/local.yml")
	if err != nil {
		panic(err)
	}

	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger = log.NewLog(conf)
	j = jwt.NewJwt(conf)
	sf = sid.NewSid()

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)

	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Username: "testuser",
		Password: "password",
		RealName: "Test User",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, nil)
	mockUserRepo.EXPECT().GetByPhone(ctx, req.Phone).Return(nil, nil)
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)
	mockTm.EXPECT().Transaction(ctx, gomock.Any()).Return(nil)

	err := userService.Register(ctx, req)

	assert.NoError(t, err)
}

func TestUserService_Register_UsernameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Username: "testuser",
		Password: "password",
		RealName: "Test User",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(&model.User{}, nil)

	err := userService.Register(ctx, req)

	assert.Error(t, err)
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	req := &v1.LoginRequest{
		UsernameOrPhone: "testuser",
		Password:        "password",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Error("failed to hash password")
	}

	mockUserRepo.EXPECT().GetByUsername(ctx, req.UsernameOrPhone).Return(&model.User{
		ID:       123,
		Username: "testuser",
		Password: string(hashedPassword),
	}, nil)

	loginResp, err := userService.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, loginResp)
	assert.NotEmpty(t, loginResp.AccessToken)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	req := &v1.LoginRequest{
		UsernameOrPhone: "testuser",
		Password:        "password",
	}

	mockUserRepo.EXPECT().GetByUsername(ctx, req.UsernameOrPhone).Return(nil, errors.New("user not found"))

	_, err := userService.Login(ctx, req)

	assert.Error(t, err)
}

func TestUserService_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	userId := int64(123)

	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(&model.User{
		ID:       userId,
		Username: "testuser",
		Email:    "test@example.com",
	}, nil)

	user, err := userService.GetProfile(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, userId, user.ID)
}

func TestUserService_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	userId := int64(123)
	req := &v1.UpdateProfileRequest{
		RealName: "testuser",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(&model.User{
		ID:    userId,
		Email: "old@example.com",
	}, nil)
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	err := userService.UpdateProfile(ctx, userId, req)

	assert.NoError(t, err)
}

func TestUserService_UpdateProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	userId := int64(123)
	req := &v1.UpdateProfileRequest{
		RealName: "testuser",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(nil, errors.New("user not found"))

	err := userService.UpdateProfile(ctx, userId, req)

	assert.Error(t, err)
}
