package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"backend/internal/model"
	"backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupParcelRepository(t *testing.T) (repository.ParcelRepository, sqlmock.Sqlmock) {
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

	repo := repository.NewRepository(logger, db)
	parcelRepo := repository.NewParcelRepository(repo)

	return parcelRepo, mock
}

func TestParcelRepository_Create(t *testing.T) {
	parcelRepo, mock := setupParcelRepository(t)

	ctx := context.Background()
	shelfID := int64(1)
	receivedAt := time.Now()
	parcel := &model.Parcel{
		ID:             1,
		TrackingNumber: "TEST123456",
		PickupCode:     "A12345",
		RecipientName:  "Test User",
		RecipientPhone: "13800138000",
		CourierCompany: "SF Express",
		Size:           model.ParcelSizeSmall,
		Weight:         1.5,
		Status:         model.ParcelStatusReceived,
		ShelfID:        &shelfID,
		ReceivedAt:     receivedAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `parcels`").
		WithArgs(
			sqlmock.AnyArg(), // tracking_number
			sqlmock.AnyArg(), // pickup_code
			sqlmock.AnyArg(), // recipient_name
			sqlmock.AnyArg(), // recipient_phone
			sqlmock.AnyArg(), // recipient_id_card
			sqlmock.AnyArg(), // courier_company
			sqlmock.AnyArg(), // size
			sqlmock.AnyArg(), // weight
			sqlmock.AnyArg(), // status
			sqlmock.AnyArg(), // shelf_id
			sqlmock.AnyArg(), // received_at
			sqlmock.AnyArg(), // shelved_at
			sqlmock.AnyArg(), // picked_up_at
			sqlmock.AnyArg(), // expected_overdue_at
			sqlmock.AnyArg(), // notes
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // id
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := parcelRepo.Create(ctx, parcel)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestParcelRepository_Update(t *testing.T) {
	parcelRepo, mock := setupParcelRepository(t)

	ctx := context.Background()
	shelfID := int64(1)
	parcel := &model.Parcel{
		ID:             1,
		TrackingNumber: "TEST123456",
		PickupCode:     "A12345",
		RecipientName:  "Test User",
		RecipientPhone: "13800138000",
		CourierCompany: "SF Express",
		Size:           model.ParcelSizeSmall,
		Weight:         1.5,
		Status:         model.ParcelStatusPickedUp,
		ShelfID:        &shelfID,
		ReceivedAt:     time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `parcels`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := parcelRepo.Update(ctx, parcel)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestParcelRepository_GetByID(t *testing.T) {
	parcelRepo, mock := setupParcelRepository(t)

	ctx := context.Background()
	parcelID := int64(1)
	shelfID := int64(1)

	// Expect the parcel query
	rows := sqlmock.NewRows([]string{
		"id", "tracking_number", "pickup_code", "recipient_name", "recipient_phone",
		"courier_company", "size", "weight", "status", "shelf_id", "received_at",
		"created_at", "updated_at",
	}).AddRow(
		parcelID, "TEST123456", "A12345", "Test User", "13800138000",
		"SF Express", "small", 1.5, "received", shelfID, time.Now(),
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `parcels`").WillReturnRows(rows)

	// Expect the shelf preload query
	shelfRows := sqlmock.NewRows([]string{
		"id", "shelf_code", "area", "floor", "column", "type", "capacity",
		"current_count", "status", "created_at", "updated_at",
	}).AddRow(
		shelfID, "A-1-1", "A", 1, 1, "small", 100, 10, "active",
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `shelves`").WillReturnRows(shelfRows)

	parcel, err := parcelRepo.GetByID(ctx, parcelID)
	assert.NoError(t, err)
	assert.NotNil(t, parcel)
	assert.Equal(t, parcelID, parcel.ID)
	assert.Equal(t, "TEST123456", parcel.TrackingNumber)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestParcelRepository_GetByPickupCode(t *testing.T) {
	parcelRepo, mock := setupParcelRepository(t)

	ctx := context.Background()
	pickupCode := "A12345"
	shelfID := int64(1)

	// Expect the parcel query
	rows := sqlmock.NewRows([]string{
		"id", "tracking_number", "pickup_code", "recipient_name", "recipient_phone",
		"courier_company", "size", "weight", "status", "shelf_id", "received_at",
		"created_at", "updated_at",
	}).AddRow(
		1, "TEST123456", pickupCode, "Test User", "13800138000",
		"SF Express", "small", 1.5, "received", shelfID, time.Now(),
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `parcels` WHERE pickup_code = \\?").WillReturnRows(rows)

	// Expect the shelf preload query
	shelfRows := sqlmock.NewRows([]string{
		"id", "shelf_code", "area", "floor", "column", "type", "capacity",
		"current_count", "status", "created_at", "updated_at",
	}).AddRow(
		shelfID, "A-1-1", "A", 1, 1, "small", 100, 10, "active",
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `shelves`").WillReturnRows(shelfRows)

	parcel, err := parcelRepo.GetByPickupCode(ctx, pickupCode)
	assert.NoError(t, err)
	assert.NotNil(t, parcel)
	assert.Equal(t, pickupCode, parcel.PickupCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestParcelRepository_GetByTrackingNumber(t *testing.T) {
	parcelRepo, mock := setupParcelRepository(t)

	ctx := context.Background()
	trackingNumber := "TEST123456"
	shelfID := int64(1)

	// Expect the parcel query
	rows := sqlmock.NewRows([]string{
		"id", "tracking_number", "pickup_code", "recipient_name", "recipient_phone",
		"courier_company", "size", "weight", "status", "shelf_id", "received_at",
		"created_at", "updated_at",
	}).AddRow(
		1, trackingNumber, "A12345", "Test User", "13800138000",
		"SF Express", "small", 1.5, "received", shelfID, time.Now(),
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `parcels` WHERE tracking_number = \\?").WillReturnRows(rows)

	// Expect the shelf preload query
	shelfRows := sqlmock.NewRows([]string{
		"id", "shelf_code", "area", "floor", "column", "type", "capacity",
		"current_count", "status", "created_at", "updated_at",
	}).AddRow(
		shelfID, "A-1-1", "A", 1, 1, "small", 100, 10, "active",
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `shelves`").WillReturnRows(shelfRows)

	parcel, err := parcelRepo.GetByTrackingNumber(ctx, trackingNumber)
	assert.NoError(t, err)
	assert.NotNil(t, parcel)
	assert.Equal(t, trackingNumber, parcel.TrackingNumber)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestParcelRepository_CountByStatus(t *testing.T) {
	parcelRepo, mock := setupParcelRepository(t)

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"status", "count"}).
		AddRow("received", 10).
		AddRow("picked_up", 5)

	mock.ExpectQuery("SELECT status, COUNT\\(\\*\\) as count FROM `parcels`").WillReturnRows(rows)

	counts, err := parcelRepo.CountByStatus(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, counts)
	assert.Equal(t, 10, counts["received"])
	assert.Equal(t, 5, counts["picked_up"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestParcelRepository_CountTodayReceived(t *testing.T) {
	parcelRepo, mock := setupParcelRepository(t)

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(15)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `parcels` WHERE received_at >= \\?").WillReturnRows(rows)

	count, err := parcelRepo.CountTodayReceived(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 15, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestParcelRepository_CountTodayPickedUp(t *testing.T) {
	parcelRepo, mock := setupParcelRepository(t)

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(8)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `parcels` WHERE \\(picked_up_at >= \\?").WillReturnRows(rows)

	count, err := parcelRepo.CountTodayPickedUp(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 8, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestParcelRepository_CountOverdue(t *testing.T) {
	parcelRepo, mock := setupParcelRepository(t)

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(3)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `parcels` WHERE status = \\?").WillReturnRows(rows)

	count, err := parcelRepo.CountOverdue(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
