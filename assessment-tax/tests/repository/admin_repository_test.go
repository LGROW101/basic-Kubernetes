package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/repository"
	"github.com/stretchr/testify/assert"
)

func TestAdminRepository_GetConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewAdminRepository(db)

	rows := sqlmock.NewRows([]string{"id", "personal_deduction", "k_receipt", "created_at", "updated_at"})
	mock.ExpectQuery("^SELECT id, personal_deduction, k_receipt, created_at, updated_at FROM admin_configs ORDER BY id DESC LIMIT 1$").WillReturnRows(rows)

	config, err := repo.GetConfig()
	assert.NoError(t, err)
	assert.Nil(t, config)

	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()
	rows = sqlmock.NewRows([]string{"id", "personal_deduction", "k_receipt", "created_at", "updated_at"}).
		AddRow(1, 60000.0, 30000.0, createdAt, updatedAt)
	mock.ExpectQuery("^SELECT id, personal_deduction, k_receipt, created_at, updated_at FROM admin_configs ORDER BY id DESC LIMIT 1$").WillReturnRows(rows)

	expectedConfig := &model.AdminConfig{
		ID:                1,
		PersonalDeduction: 60000.0,
		KReceipt:          30000.0,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	config, err = repo.GetConfig()
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)
}

func TestAdminRepository_UpdateConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewAdminRepository(db)

	config := &model.AdminConfig{
		PersonalDeduction: 70000.0,
		KReceipt:          40000.0,
	}

	mock.ExpectExec("^UPDATE admin_configs SET personal_deduction = \\$1, k_receipt = \\$2, updated_at = NOW\\(\\) WHERE id = 1$").
		WithArgs(config.PersonalDeduction, config.KReceipt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateConfig(config)
	assert.NoError(t, err)
}
func TestAdminRepository_InsertConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewAdminRepository(db)

	config := &model.AdminConfig{
		PersonalDeduction: 60000.0,
		KReceipt:          30000.0,
	}

	mock.ExpectExec("^INSERT INTO admin_configs \\(personal_deduction, k_receipt\\) VALUES \\(\\$1, \\$2\\)$").
		WithArgs(config.PersonalDeduction, config.KReceipt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.InsertConfig(config)
	assert.NoError(t, err)
}
