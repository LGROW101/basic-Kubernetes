package service_test

import (
	"testing"

	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/service"
	"github.com/LGROW101/assessment-tax/tests/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaxCalculatorService_GetAllCalculations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaxRepository(ctrl)
	expectedCalculations := []*model.TaxCalculation{
		{
			TotalIncome:       1000000.0,
			WHT:               100000.0,
			PersonalAllowance: 60000.0,
			Donation:          10000.0,
			KReceipt:          20000.0,
			Tax:               89500.0,
			TaxPayable:        89500.0,
			TaxLevel: []model.TaxRate{
				{Level: "0-150,000", Tax: 0},
				{Level: "150,001-500,000", Tax: 0},
				{Level: "500,001-1,000,000", Tax: 89500.0},
				{Level: "1,000,001-2,000,000", Tax: 0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0},
			},
		},
	}
	mockRepo.EXPECT().GetAllCalculations().Return(expectedCalculations, nil)

	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	taxSvc := service.NewTaxCalculatorService(mockRepo, mockAdminRepo)
	calculations, err := taxSvc.GetAllCalculations()

	assert.NoError(t, err)
	assert.Equal(t, expectedCalculations, calculations)
}
func TestTaxCalculatorService_CalculateTax(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var expectedTaxCalculation *model.TaxCalculation

	mockRepo := mocks.NewMockTaxRepository(ctrl)
	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)

	config := &model.AdminConfig{
		PersonalDeduction: 60000,
		KReceipt:          30000,
	}
	mockAdminRepo.EXPECT().GetConfig().Return(config, nil)

	iotalIncome := 1000000.0
	wht := 100000.0
	allowances := []model.Allowance{
		{AllowanceType: "donation", Amount: 10000.0},
		{AllowanceType: "k-receipt", Amount: 20000.0},
	}

	taxableIncome := iotalIncome - config.PersonalDeduction - allowances[0].Amount - allowances[1].Amount
	var tax float64
	switch {
	case taxableIncome <= 150000:
		tax = taxableIncome * 0.05
	case taxableIncome <= 500000:
		tax = 7500 + (taxableIncome-150000)*0.10
	case taxableIncome <= 1000000:
		tax = 35000 + (taxableIncome-500000)*0.15
	case taxableIncome <= 2000000:
		tax = 110000 + (taxableIncome-1000000)*0.20
	default:
		tax = 310000 + (taxableIncome-2000000)*0.35
	}

	taxPayable := tax - wht
	if taxPayable < 0 {
		taxPayable = 0
	}
	taxRefund := wht - tax
	if taxRefund < 0 {
		taxRefund = 0
	}

	expectedTaxLevel := []model.TaxRate{
		{Level: "0-150,000", Tax: 0},
		{Level: "150,001-500,000", Tax: 0},
		{Level: "500,001-1,000,000", Tax: 0},
		{Level: "1,000,001-2,000,000", Tax: 0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0},
	}
	switch {
	case taxableIncome <= 150000:
		expectedTaxLevel[0].Tax = taxPayable
	case taxableIncome <= 500000:
		expectedTaxLevel[1].Tax = taxPayable
	case taxableIncome <= 1000000:
		expectedTaxLevel[2].Tax = taxPayable
	case taxableIncome <= 2000000:
		expectedTaxLevel[3].Tax = taxPayable
	default:
		expectedTaxLevel[4].Tax = taxPayable
	}

	// Initialize expectedTaxCalculation with the expected values
	expectedTaxCalculation = &model.TaxCalculation{
		TotalIncome:       iotalIncome,
		WHT:               wht,
		PersonalAllowance: config.PersonalDeduction,
		Donation:          allowances[0].Amount,
		KReceipt:          allowances[1].Amount,
		Tax:               tax,
		TaxPayable:        taxPayable,
		TaxRefund:         taxRefund,
		TaxLevel:          expectedTaxLevel,
	}

	mockRepo.EXPECT().Save(expectedTaxCalculation).Return(nil)
	expectedTaxCalculationResponse := &model.TaxCalculationResponse{
		Tax:       nil,
		TaxRefund: &taxRefund,
		TaxLevel:  expectedTaxLevel,
	}

	taxSvc := service.NewTaxCalculatorService(mockRepo, mockAdminRepo)
	taxCalculation, err := taxSvc.CalculateTax(iotalIncome, wht, allowances)

	assert.NoError(t, err)
	assert.Equal(t, expectedTaxCalculationResponse, taxCalculation)
}
