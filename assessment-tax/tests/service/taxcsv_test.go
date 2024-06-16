package service_test

import (
	"strings"
	"testing"

	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/service"
	"github.com/LGROW101/assessment-tax/tests/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaxCSVService_ImportCSV(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taxRepo := mocks.NewMockTaxRepository(ctrl)
	adminRepo := mocks.NewMockAdminRepository(ctrl)

	adminConfig := &model.AdminConfig{
		PersonalDeduction: 60000,
	}
	adminRepo.EXPECT().GetConfig().Return(adminConfig, nil).Times(3)

	csvData := `income,wht,donation
   500000,0,0
   600000,40000,20000
   750000,50000,15000`

	expectedResult := []map[string]float64{
		{"totalIncome": 500000, "tax": 29000},
		{"totalIncome": 600000, "taxRefund": 2000},
		{"totalIncome": 750000, "tax": 11250},
	}

	taxCSVService := service.NewTaxCSVService(taxRepo, adminRepo)
	result, err := taxCSVService.ImportCSV(strings.NewReader(csvData))

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestTaxCSVService_CalculateTax(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	taxRepo := mocks.NewMockTaxRepository(ctrl)
	adminRepo := mocks.NewMockAdminRepository(ctrl)
	adminConfig := &model.AdminConfig{
		PersonalDeduction: 60000,
	}
	adminRepo.EXPECT().GetConfig().Return(adminConfig, nil).Times(3)
	taxCSVService := service.NewTaxCSVService(taxRepo, adminRepo)

	testCases := []struct {
		name     string
		income   float64
		wht      float64
		donation float64
		expected float64
	}{
		{"Case 1", 500000, 0, 0, 29000},
		{"Case 2", 600000, 40000, 20000, 0},
		{"Case 3", 750000, 50000, 15000, 11250},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _, err := taxCSVService.CalculateTax(tc.income, tc.wht, tc.donation)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseFields(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []float64
		hasError bool
	}{
		{"Valid input", []string{"500000", "50000", "0"}, []float64{500000, 50000, 0}, false},
		{"Missing fields", []string{"600000", "40000", "20000"}, []float64{600000, 40000, 20000}, false},
		{"Invalid income", []string{"invalid", "50000", "0"}, nil, true},
		{"Invalid wht", []string{"500000", "invalid", "0"}, nil, true},
		{"Invalid donation", []string{"500000", "50000", "invalid"}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			income, wht, donation, err := service.ParseFields(tc.input)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected[0], income)
				assert.Equal(t, tc.expected[1], wht)
				assert.Equal(t, tc.expected[2], donation)
			}
		})
	}
}

func TestParseDonation(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"Valid amount", "10000", 10000, false},
		{"Valid percentage", "5%", 0.05, false},
		{"Invalid input", "invalid", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.ParseDonation(tc.input)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
