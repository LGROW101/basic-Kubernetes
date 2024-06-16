package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LGROW101/assessment-tax/handler"
	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/tests/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTax(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaxCalculatorService(ctrl)
	calculatorHandler := handler.NewCalculatorHandler(mockService)

	totalIncome := 1000000.0
	wht := 50000.0
	allowances := []model.Allowance{
		{AllowanceType: "allowance1", Amount: 10000.0},
		{AllowanceType: "allowance2", Amount: 20000.0},
	}
	includeTaxLevel := true

	expectedTax := 100000.0
	expectedTaxResponse := &model.TaxCalculationResponse{
		Tax: &expectedTax,
		TaxLevel: []model.TaxRate{
			{Level: "Level 2", Tax: 10000.0},
		},
	}

	mockService.EXPECT().CalculateTax(totalIncome, wht, allowances).Return(expectedTaxResponse, nil)

	reqBody, _ := json.Marshal(handler.CalculateTaxRequest{
		TotalIncome:     totalIncome,
		WHT:             wht,
		Allowances:      allowances,
		IncludeTaxLevel: includeTaxLevel,
	})

	req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := calculatorHandler.CalculateTax(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTax, response["tax"].(float64))

	taxLevelResponse, ok := response["taxLevel"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, taxLevelResponse, 1)

	taxRate := taxLevelResponse[0].(map[string]interface{})
	assert.Equal(t, expectedTaxResponse.TaxLevel[0].Level, taxRate["level"])
	assert.Equal(t, expectedTaxResponse.TaxLevel[0].Tax, taxRate["tax"])
}

func TestCalculateTaxWithInvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaxCalculatorService(ctrl)
	calculatorHandler := handler.NewCalculatorHandler(mockService)

	// Add this line to set an empty expectation for CalculateTax
	mockService.EXPECT().CalculateTax(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	invalidReqBody := []byte(`{"invalidField": "invalid"}`)

	req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(invalidReqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := calculatorHandler.CalculateTax(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestGetAllCalculations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaxCalculatorService(ctrl)
	calculatorHandler := handler.NewCalculatorHandler(mockService)

	expectedCalculations := []*model.TaxCalculation{
		{
			TaxPayable: 100000.0,
			TaxLevel: []model.TaxRate{
				{Level: "Level 1", Tax: 5000.0},
			},
		},
		{
			TaxPayable: 200000.0,
			TaxLevel: []model.TaxRate{
				{Level: "Level 2", Tax: 10000.0},
			},
		},
	}

	mockService.EXPECT().GetAllCalculations().Return(expectedCalculations, nil)

	req := httptest.NewRequest(http.MethodGet, "/calculations", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := calculatorHandler.GetAllCalculations(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []*model.TaxCalculation
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCalculations, response)
}
func TestCalculateTaxWithoutTaxLevel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaxCalculatorService(ctrl)
	calculatorHandler := handler.NewCalculatorHandler(mockService)

	totalIncome := 1000000.0
	wht := 50000.0
	allowances := []model.Allowance{
		{AllowanceType: "allowance1", Amount: 10000.0},
		{AllowanceType: "allowance2", Amount: 20000.0},
	}
	includeTaxLevel := false

	expectedTax := 100000.0
	expectedTaxResponse := &model.TaxCalculationResponse{
		Tax: &expectedTax,
	}

	mockService.EXPECT().CalculateTax(totalIncome, wht, allowances).Return(expectedTaxResponse, nil)

	reqBody, _ := json.Marshal(handler.CalculateTaxRequest{
		TotalIncome:     totalIncome,
		WHT:             wht,
		Allowances:      allowances,
		IncludeTaxLevel: includeTaxLevel,
	})

	req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := calculatorHandler.CalculateTax(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTax, response["tax"].(float64))
	assert.Nil(t, response["taxLevel"])
}
func TestGetAllCalculationsWithServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaxCalculatorService(ctrl)
	calculatorHandler := handler.NewCalculatorHandler(mockService)

	mockService.EXPECT().GetAllCalculations().Return(nil, errors.New("service error"))

	req := httptest.NewRequest(http.MethodGet, "/calculations", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := calculatorHandler.GetAllCalculations(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestCalculateTaxWithInvalidRequestValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaxCalculatorService(ctrl)
	calculatorHandler := handler.NewCalculatorHandler(mockService)

	invalidReqBodies := []handler.CalculateTaxRequest{
		{TotalIncome: -1000, WHT: 50000, Allowances: []model.Allowance{{AllowanceType: "allowance1", Amount: 10000}}},
		{TotalIncome: 1000000, WHT: -50000, Allowances: []model.Allowance{{AllowanceType: "allowance1", Amount: 10000}}},
		{TotalIncome: 1000000, WHT: 50000, Allowances: []model.Allowance{}},
	}

	for _, reqBody := range invalidReqBodies {
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		err := calculatorHandler.CalculateTax(c)
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	}
}

func TestCalculateTaxWithServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaxCalculatorService(ctrl)
	calculatorHandler := handler.NewCalculatorHandler(mockService)

	totalIncome := 1000000.0
	wht := 50000.0
	allowances := []model.Allowance{
		{AllowanceType: "allowance1", Amount: 10000.0},
	}

	mockService.EXPECT().CalculateTax(totalIncome, wht, allowances).Return(nil, errors.New("service error"))

	reqBody, _ := json.Marshal(handler.CalculateTaxRequest{
		TotalIncome: totalIncome,
		WHT:         wht,
		Allowances:  allowances,
	})

	req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := calculatorHandler.CalculateTax(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}
