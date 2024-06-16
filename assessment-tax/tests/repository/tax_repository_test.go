package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/repository"
	"github.com/stretchr/testify/assert"
)

func TestTaxRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewTaxRepository(db)

	taxCalculation := &model.TaxCalculation{
		TotalIncome:       1000000,
		WHT:               100000,
		PersonalAllowance: 60000,
		Donation:          10000,
		KReceipt:          30000,
		Tax:               200000,
	}

	mock.ExpectExec("^INSERT INTO tax_calculations").
		WithArgs(taxCalculation.TotalIncome, taxCalculation.WHT, taxCalculation.PersonalAllowance, taxCalculation.Donation, taxCalculation.KReceipt, taxCalculation.Tax).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Save(taxCalculation)
	assert.NoError(t, err)

	mock.ExpectExec("^INSERT INTO tax_calculations").
		WithArgs(taxCalculation.TotalIncome, taxCalculation.WHT, taxCalculation.PersonalAllowance, taxCalculation.Donation, taxCalculation.KReceipt, taxCalculation.Tax).
		WillReturnError(errors.New("database error"))

	err = repo.Save(taxCalculation)
	assert.Error(t, err)
}

func TestTaxRepository_GetAllCalculations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewTaxRepository(db)

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "totalIncome", "wht", "personal_allowance", "donation", "k_receipt", "tax", "created_at"}).
		AddRow(1, 1000000, 100000, 60000, 10000, 30000, 200000, createdAt).
		AddRow(2, 800000, 80000, 60000, 5000, 20000, 150000, createdAt)

	mock.ExpectQuery("^SELECT id, totalIncome, wht, personal_allowance, donation, k_receipt, tax, created_at FROM tax_calculations$").
		WillReturnRows(rows)

	expectedCalculations := []*model.TaxCalculation{
		{
			ID:                1,
			TotalIncome:       1000000,
			WHT:               100000,
			PersonalAllowance: 60000,
			Donation:          10000,
			KReceipt:          30000,
			Tax:               200000,
			CreatedAt:         createdAt,
		},
		{
			ID:                2,
			TotalIncome:       800000,
			WHT:               80000,
			PersonalAllowance: 60000,
			Donation:          5000,
			KReceipt:          20000,
			Tax:               150000,
			CreatedAt:         createdAt,
		},
	}

	calculations, err := repo.GetAllCalculations()
	assert.NoError(t, err)
	assert.Equal(t, expectedCalculations, calculations)

	mock.ExpectQuery("^SELECT id, totalIncome, wht, personal_allowance, donation, k_receipt, tax, created_at FROM tax_calculations$").
		WillReturnError(errors.New("database error"))

	calculations, err = repo.GetAllCalculations()
	assert.Error(t, err)
	assert.Nil(t, calculations)

	rows = sqlmock.NewRows([]string{"id", "totalIncome", "wht", "personal_allowance", "donation", "k_receipt", "tax", "created_at"}).
		AddRow(1, "invalid", 100000, 60000, 10000, 30000, 200000, createdAt)

	mock.ExpectQuery("^SELECT id, totalIncome, wht, personal_allowance, donation, k_receipt, tax, created_at FROM tax_calculations$").
		WillReturnRows(rows)

	calculations, err = repo.GetAllCalculations()
	assert.Error(t, err)
	assert.Nil(t, calculations)
}
