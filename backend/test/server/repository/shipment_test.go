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

func setupShipmentRepository(t *testing.T) (repository.ShipmentRepository, sqlmock.Sqlmock) {
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
	shipmentRepo := repository.NewShipmentRepository(repo)

	return shipmentRepo, mock
}

func TestShipmentRepository_Create(t *testing.T) {
	shipmentRepo, mock := setupShipmentRepository(t)

	ctx := context.Background()
	shipment := &model.Shipment{
		ID:               1,
		ShipmentNumber:   "SH202412345678",
		SenderID:         1,
		SenderName:       "Test Sender",
		SenderPhone:      "13800138000",
		RecipientName:    "Test Recipient",
		RecipientPhone:   "13900139000",
		RecipientAddress: "Test Address",
		CourierCompany:   "SF Express",
		Weight:           1.5,
		Volume:           0.01,
		Freight:          15.0,
		Status:           model.ShipmentStatusPending,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `shipments`").
		WithArgs(
			sqlmock.AnyArg(), // shipment_number
			sqlmock.AnyArg(), // sender_id
			sqlmock.AnyArg(), // sender_name
			sqlmock.AnyArg(), // sender_phone
			sqlmock.AnyArg(), // recipient_name
			sqlmock.AnyArg(), // recipient_phone
			sqlmock.AnyArg(), // recipient_address
			sqlmock.AnyArg(), // courier_company
			sqlmock.AnyArg(), // weight
			sqlmock.AnyArg(), // volume
			sqlmock.AnyArg(), // freight
			sqlmock.AnyArg(), // status
			sqlmock.AnyArg(), // shipped_at
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // id
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := shipmentRepo.Create(ctx, shipment)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestShipmentRepository_Update(t *testing.T) {
	shipmentRepo, mock := setupShipmentRepository(t)

	ctx := context.Background()
	shipment := &model.Shipment{
		ID:               1,
		ShipmentNumber:   "SH202412345678",
		SenderID:         1,
		SenderName:       "Test Sender",
		SenderPhone:      "13800138000",
		RecipientName:    "Test Recipient",
		RecipientPhone:   "13900139000",
		RecipientAddress: "Test Address",
		CourierCompany:   "SF Express",
		Weight:           1.5,
		Volume:           0.01,
		Freight:          15.0,
		Status:           model.ShipmentStatusShipped,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `shipments`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := shipmentRepo.Update(ctx, shipment)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestShipmentRepository_GetByID(t *testing.T) {
	shipmentRepo, mock := setupShipmentRepository(t)

	ctx := context.Background()
	shipmentID := int64(1)

	rows := sqlmock.NewRows([]string{
		"id", "shipment_number", "sender_id", "sender_name", "sender_phone",
		"recipient_name", "recipient_phone", "recipient_address", "courier_company",
		"weight", "volume", "freight", "status", "created_at", "updated_at",
	}).AddRow(
		shipmentID, "SH202412345678", 1, "Test Sender", "13800138000",
		"Test Recipient", "13900139000", "Test Address", "SF Express",
		1.5, 0.01, 15.0, "pending", time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `shipments`").WillReturnRows(rows)

	shipment, err := shipmentRepo.GetByID(ctx, shipmentID)
	assert.NoError(t, err)
	assert.NotNil(t, shipment)
	assert.Equal(t, shipmentID, shipment.ID)
	assert.Equal(t, "SH202412345678", shipment.ShipmentNumber)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestShipmentRepository_GetByShipmentNumber(t *testing.T) {
	shipmentRepo, mock := setupShipmentRepository(t)

	ctx := context.Background()
	shipmentNumber := "SH202412345678"

	rows := sqlmock.NewRows([]string{
		"id", "shipment_number", "sender_id", "sender_name", "sender_phone",
		"recipient_name", "recipient_phone", "recipient_address", "courier_company",
		"weight", "volume", "freight", "status", "created_at", "updated_at",
	}).AddRow(
		1, shipmentNumber, 1, "Test Sender", "13800138000",
		"Test Recipient", "13900139000", "Test Address", "SF Express",
		1.5, 0.01, 15.0, "pending", time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `shipments` WHERE shipment_number = \\?").WillReturnRows(rows)

	shipment, err := shipmentRepo.GetByShipmentNumber(ctx, shipmentNumber)
	assert.NoError(t, err)
	assert.NotNil(t, shipment)
	assert.Equal(t, shipmentNumber, shipment.ShipmentNumber)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestShipmentRepository_ListBySenderID(t *testing.T) {
	shipmentRepo, mock := setupShipmentRepository(t)

	ctx := context.Background()
	senderID := int64(1)
	status := "pending"
	page := 1
	pageSize := 20

	// Count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `shipments` WHERE sender_id = \\? AND status = \\?").
		WithArgs(senderID, status).
		WillReturnRows(countRows)

	// List query - include limit parameter
	rows := sqlmock.NewRows([]string{
		"id", "shipment_number", "sender_id", "sender_name", "sender_phone",
		"recipient_name", "recipient_phone", "recipient_address", "courier_company",
		"weight", "volume", "freight", "status", "created_at", "updated_at",
	}).AddRow(
		1, "SH202412345678", senderID, "Test Sender", "13800138000",
		"Test Recipient", "13900139000", "Test Address", "SF Express",
		1.5, 0.01, 15.0, status, time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `shipments` WHERE sender_id = \\? AND status = \\?").
		WithArgs(senderID, status, pageSize).
		WillReturnRows(rows)

	shipments, total, err := shipmentRepo.ListBySenderID(ctx, senderID, status, page, pageSize)
	assert.NoError(t, err)
	assert.NotNil(t, shipments)
	assert.Equal(t, 1, len(shipments))
	assert.Equal(t, int64(1), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}
