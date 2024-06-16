package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LGROW101/assessment-tax/handler"
	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/tests/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	adminHandler := handler.NewAdminHandler(mockAdminRepo)

	expectedConfig := &model.AdminConfig{
		PersonalDeduction: 60000,
		KReceipt:          30000,
	}

	mockAdminRepo.EXPECT().GetConfig().Return(expectedConfig, nil)

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := adminHandler.GetConfig(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.AdminConfig
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, &response)
}

func TestGetConfigWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	adminHandler := handler.NewAdminHandler(mockAdminRepo)

	mockAdminRepo.EXPECT().GetConfig().Return(nil, errors.New("repository error"))

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := adminHandler.GetConfig(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestUpdateConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	adminHandler := handler.NewAdminHandler(mockAdminRepo)

	existingConfig := &model.AdminConfig{
		PersonalDeduction: 60000,
		KReceipt:          30000,
	}

	mockAdminRepo.EXPECT().GetConfig().Return(existingConfig, nil)
	mockAdminRepo.EXPECT().UpdateConfig(gomock.Any()).DoAndReturn(func(config *model.AdminConfig) error {
		assert.Equal(t, float64(70000), config.PersonalDeduction)
		assert.Equal(t, existingConfig.KReceipt, config.KReceipt)
		return nil
	})

	reqBody := `{"personalDeduction":70000}`
	req := httptest.NewRequest(http.MethodPut, "/config", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := adminHandler.UpdateConfig(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.AdminResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(70000), response.PersonalDeduction)
	assert.Equal(t, float64(0), response.KReceipt)
}

func TestUpdateConfigWithInvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	adminHandler := handler.NewAdminHandler(mockAdminRepo)

	reqBody := `{"personalDeduction":"invalid","kReceipt":40000}`
	req := httptest.NewRequest(http.MethodPut, "/config", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := adminHandler.UpdateConfig(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestUpdateConfigWithValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	adminHandler := handler.NewAdminHandler(mockAdminRepo)

	originalConfig := &model.AdminConfig{
		PersonalDeduction: 60000,
		KReceipt:          30000,
	}

	reqBody := `{"personalDeduction":-1,"kReceipt":40000}`
	req := httptest.NewRequest(http.MethodPut, "/config", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockAdminRepo.EXPECT().GetConfig().Return(originalConfig, nil)

	err := adminHandler.UpdateConfig(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}
func TestGetConfigNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	adminHandler := handler.NewAdminHandler(mockAdminRepo)

	mockAdminRepo.EXPECT().GetConfig().Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := adminHandler.GetConfig(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)
}

func TestUpdateConfigInsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	adminHandler := handler.NewAdminHandler(mockAdminRepo)

	mockAdminRepo.EXPECT().GetConfig().Return(nil, nil)
	mockAdminRepo.EXPECT().InsertConfig(gomock.Any()).Return(errors.New("insert error"))

	reqBody := `{"personalDeduction":70000,"kReceipt":40000}`
	req := httptest.NewRequest(http.MethodPut, "/config", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := adminHandler.UpdateConfig(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestUpdateConfigUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	adminHandler := handler.NewAdminHandler(mockAdminRepo)

	existingConfig := &model.AdminConfig{
		PersonalDeduction: 60000,
		KReceipt:          30000,
	}

	mockAdminRepo.EXPECT().GetConfig().Return(existingConfig, nil)
	mockAdminRepo.EXPECT().UpdateConfig(gomock.Any()).Return(errors.New("update error"))

	reqBody := `{"personalDeduction":70000}`
	req := httptest.NewRequest(http.MethodPut, "/config", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := adminHandler.UpdateConfig(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}
