// csv_test
package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LGROW101/assessment-tax/handler"
	"github.com/LGROW101/assessment-tax/tests/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUploadCSV(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVService := mocks.NewMockTaxCSVService(ctrl)
	csvHandler := handler.NewCSVHandler(mockCSVService)

	// Create a new HTTP request with multipart form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxFile", "test.csv")
	assert.NoError(t, err)
	part.Write([]byte("csv data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	expectedTaxes := []map[string]float64{
		{"totalIncome": 500000, "tax": 50000},
		{"totalIncome": 1000000, "tax": 200000},
	}
	mockCSVService.EXPECT().ImportCSV(gomock.Any()).Return(expectedTaxes, nil)

	err = csvHandler.UploadCSV(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	taxes, ok := response["taxes"].([]interface{})
	assert.True(t, ok)

	for i, tax := range taxes {
		taxMap, ok := tax.(map[string]interface{})
		assert.True(t, ok)

		expectedTax := expectedTaxes[i]
		assertEqualFloat64Maps(t, expectedTax, convertMapToFloat64(taxMap))
	}
}

func TestUploadCSVWithMissingFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVService := mocks.NewMockTaxCSVService(ctrl)
	csvHandler := handler.NewCSVHandler(mockCSVService)

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := csvHandler.UploadCSV(c)
	assert.Error(t, err)
	assert.IsType(t, &echo.HTTPError{}, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestUploadCSVWithReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVService := mocks.NewMockTaxCSVService(ctrl)
	csvHandler := handler.NewCSVHandler(mockCSVService)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxFile", "test.csv")
	assert.NoError(t, err)
	part.Write([]byte("csv data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockCSVService.EXPECT().ImportCSV(gomock.Any()).Return(nil, errors.New("read error"))

	err = csvHandler.UploadCSV(c)
	assert.Error(t, err)
	assert.IsType(t, &echo.HTTPError{}, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func convertMapToFloat64(m map[string]interface{}) map[string]float64 {
	result := make(map[string]float64)
	for k, v := range m {
		result[k], _ = v.(float64)
	}
	return result
}

func assertEqualFloat64Maps(t *testing.T, expected, actual map[string]float64) {
	for k, v := range expected {
		actualValue, ok := actual[k]
		if !ok {
			t.Errorf("Key %q not found in actual map", k)
			return
		}
		if actualValue != v {
			t.Errorf("Values not equal for key %q: expected %f, got %f", k, v, actualValue)
			return
		}
	}
}

func TestUploadCSVSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVService := mocks.NewMockTaxCSVService(ctrl)
	csvHandler := handler.NewCSVHandler(mockCSVService)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxFile", "test.csv")
	assert.NoError(t, err)
	part.Write([]byte("csv data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	expectedTaxes := []map[string]float64{
		{"totalIncome": 500000, "tax": 50000},
		{"totalIncome": 1000000, "tax": 200000},
	}
	mockCSVService.EXPECT().ImportCSV(gomock.Any()).Return(expectedTaxes, nil)

	err = csvHandler.UploadCSV(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	taxes, ok := response["taxes"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, len(expectedTaxes), len(taxes))

	for i, tax := range taxes {
		taxMap, ok := tax.(map[string]interface{})
		assert.True(t, ok)

		assert.Equal(t, expectedTaxes[i]["totalIncome"], taxMap["totalIncome"])
		assert.Equal(t, expectedTaxes[i]["tax"], taxMap["tax"])
	}
}

func TestUploadCSVWithoutFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVService := mocks.NewMockTaxCSVService(ctrl)
	csvHandler := handler.NewCSVHandler(mockCSVService)

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := csvHandler.UploadCSV(c)
	assert.Error(t, err)
	assert.IsType(t, &echo.HTTPError{}, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestUploadCSVServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVService := mocks.NewMockTaxCSVService(ctrl)
	csvHandler := handler.NewCSVHandler(mockCSVService)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxFile", "test.csv")
	assert.NoError(t, err)
	part.Write([]byte("csv data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockCSVService.EXPECT().ImportCSV(gomock.Any()).Return(nil, errors.New("service error"))

	err = csvHandler.UploadCSV(c)
	assert.Error(t, err)
	assert.IsType(t, &echo.HTTPError{}, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}
