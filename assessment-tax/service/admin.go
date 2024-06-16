// service/admin.go
package service

import (
	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/repository"
)

type AdminServiceInterface interface {
	GetConfig() (*model.AdminConfig, error)
	UpdateConfig(config *model.AdminConfig) error
}

type AdminService struct {
	adminRepo repository.AdminRepository
}

func NewAdminService(adminRepo repository.AdminRepository) AdminServiceInterface {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

func (s *AdminService) GetConfig() (*model.AdminConfig, error) {
	return s.adminRepo.GetConfig()
}

func (s *AdminService) UpdateConfig(config *model.AdminConfig) error {
	return s.adminRepo.UpdateConfig(config)
}
