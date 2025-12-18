package repository

import (
	"context"
	"backend/pkg/log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"backend/internal/model"
	"backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logger *log.Logger

func setupRepository(t *testing.T) (repository.UserRepository, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	//rdb, _ := redismock.NewClientMock()

	repo := repository.NewRepository(logger, db)
	userRepo := repository.NewUserRepository(repo)

	return userRepo, mock
}

func TestUserRepository_Create(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	user := &model.User{
		ID:        1,
		Username:  "testuser",
		RealName:  "Test User",
		Password:  "password",
		Phone:     "13800138000",
		Email:     "test@example.com",
		Role:      "user", // Add role field to match GORM default
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(user.Username, user.Password, user.RealName, user.Phone, user.IDCard, user.Email, user.Role, user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	user := &model.User{
		ID:        1,
		Username:  "testuser",
		RealName:  "Test User",
		Password:  "password",
		Phone:     "13800138000",
		Email:     "test@example.com",
		Role:      "user", // Add role field
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.Update(ctx, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetById(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	userId := int64(123)

	rows := sqlmock.NewRows([]string{"id", "username", "real_name", "password", "phone", "email", "role", "created_at", "updated_at"}).
		AddRow(123, "testuser", "Test User", "password", "13800138000", "test@example.com", "user", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(rows)

	user, err := userRepo.GetByID(ctx, userId)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(123), user.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	email := "test@example.com"

	rows := sqlmock.NewRows([]string{"id", "username", "real_name", "password", "phone", "email", "role", "created_at", "updated_at"}).
		AddRow(123, "testuser", "Test User", "password", "13800138000", "test@example.com", "user", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(rows)

	user, err := userRepo.GetByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}
