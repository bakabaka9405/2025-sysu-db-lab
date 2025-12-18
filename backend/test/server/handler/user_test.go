package handler

import (
	v1 "backend/api/v1"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/test/mocks/service"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RegisterRequest{
		Username: "testuser",
		Password: "123456",
		RealName: "Test User",
		Phone:    "13800138000",
		Email:    "xxx@gmail.com",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Register(gomock.Any(), &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.POST("/register", userHandler.Register)

	e := newHttpExcept(t, router)
	obj := e.POST("/register").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.LoginRequest{
		UsernameOrPhone: "testuser",
		Password:        "123456",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Login(gomock.Any(), &params).Return(&v1.LoginResponseData{
		AccessToken: "test-token",
		User: v1.UserInfo{
			ID:       123,
			Username: "testuser",
		},
	}, nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.POST("/login", userHandler.Login)

	obj := newHttpExcept(t, router).POST("/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	obj.Value("data").Object().Value("accessToken").IsEqual("test-token")
}

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testuser"
	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().GetProfile(gomock.Any(), int64(123)).Return(&v1.GetProfileResponseData{
		ID:       123,
		Username: username,
	}, nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.Use(middleware.NoStrictAuth(jwt, logger))
	router.GET("/user", userHandler.GetProfile)

	obj := newHttpExcept(t, router).GET("/user").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("id").IsEqual(float64(123))
	objData.Value("username").IsEqual(username)
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.UpdateProfileRequest{
		RealName: "Alan Zhang",
		Email:    "alan@gmail.com",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().UpdateProfile(gomock.Any(), int64(123), &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.PUT("/user", userHandler.UpdateProfile)

	obj := newHttpExcept(t, router).PUT("/user").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}
