package repositories

import (
	"context"
	"errors"
	"flea-market/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *ItemRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %s", err)
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %s", err)
	}
	repo := NewItemRepository(gdb)
	return gdb, mock, repo
}

func TestItemRepository_Create_Success(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	defer mock.ExpectClose()

	item := models.Item{
		UserID:      1,
		Name:        "Item 1",
		Price:       100,
		Description: "",
		SoldOut:     false,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "items"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			item.Name,
			item.Price,
			item.Description,
			item.SoldOut,
			item.UserID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	ctx := context.Background()
	result, err := repo.Create(ctx, item)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.UserID != item.UserID || result.Name != item.Name {
		t.Errorf("item not correctly returned")
	}
}

func TestItemRepository_Create_Timeout(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	defer mock.ExpectClose()

	item := models.Item{UserID: 1, Name: "timeout", Price: 100}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "items"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			item.Name,
			item.Price,
			item.Description,
			item.SoldOut,
			item.UserID,
		).
		WillReturnError(errors.New("context deadline exceeded"))
	mock.ExpectRollback()

	ctx := context.Background()
	_, err := repo.Create(ctx, item)
	if err == nil {
		t.Errorf("expected timeout error, got %v", err)
	}

	assert.ErrorContains(t, err, "DB処理がタイムアウトしました:")
}

func TestItemRepository_Create_Canceled(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	defer mock.ExpectClose()

	item := models.Item{UserID: 1, Name: "cancel", Price: 100}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "items"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			item.Name,
			item.Price,
			item.Description,
			item.SoldOut,
			item.UserID,
		).WillReturnError(errors.New("context canceled"))
	mock.ExpectRollback()

	ctx := context.Background()
	_, err := repo.Create(ctx, item)
	if err == nil {
		t.Errorf("expected canceled error, got %v", err)
	}

	assert.ErrorContains(t, err, "DB処理がキャンセルされました:")
}

func TestItemRepository_Create_OtherError(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	defer mock.ExpectClose()

	item := models.Item{UserID: 1, Name: "othererr", Price: 100}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "items"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			item.Name,
			item.Price,
			item.Description,
			item.SoldOut,
			item.UserID,
		).WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	ctx := context.Background()
	_, err := repo.Create(ctx, item)
	if err == nil {
		t.Errorf("expected custom db error, got %v", err)
	}

	assert.ErrorContains(t, err, "duplicated key not allowed")
}

func TestItemRepository_FindById_Success(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	defer mock.ExpectClose()

	itemID := uint(1)
	userID := uint(1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items" WHERE (id = $1 AND user_id = $2) AND "items"."deleted_at" IS NULL ORDER BY "items"."id" LIMIT $3`)).
		WithArgs(itemID, userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "price"}).
			AddRow(itemID, userID, "Test", 100))

	ctx := context.Background()
	item, err := repo.FindById(ctx, itemID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if item.ID != itemID || item.UserID != userID {
		t.Errorf("item not correctly returned")
	}
}

func TestItemRepository_FindById_NotFound(t *testing.T) {
	_, mock, repo := setupTestDB(t)
	defer mock.ExpectClose()

	itemID := uint(999)
	userID := uint(1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items" WHERE (id = $1 AND user_id = $2) AND "items"."deleted_at" IS NULL ORDER BY "items"."id" LIMIT $3`)).
		WithArgs(itemID, userID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	ctx := context.Background()
	_, err := repo.FindById(ctx, itemID, userID)
	if err == nil {
		t.Fatalf("expected not found error")
	}

	assert.ErrorContains(t, err, "Not Found From DB")
}
