package service_test

import (
	"testing"

	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/service"
	"github.com/LGROW101/assessment-tax/tests/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAdminService_GetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAdminRepository(ctrl)
	expectedConfig := &model.AdminConfig{
		ID:                1,
		PersonalDeduction: 60000,
		KReceipt:          30000,
	}
	mockRepo.EXPECT().GetConfig().Return(expectedConfig, nil)

	adminSvc := service.NewAdminService(mockRepo)
	config, err := adminSvc.GetConfig()

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)
}

func TestAdminService_UpdateConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAdminRepository(ctrl)
	config := &model.AdminConfig{
		PersonalDeduction: 70000,
		KReceipt:          40000,
	}
	mockRepo.EXPECT().UpdateConfig(config).Return(nil)

	adminSvc := service.NewAdminService(mockRepo)
	err := adminSvc.UpdateConfig(config)

	assert.NoError(t, err)
}
