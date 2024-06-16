package service

import (
	"math"

	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/repository"
)

type TaxCalculatorService interface {
	GetAllCalculations() ([]*model.TaxCalculation, error)
	CalculateTax(totalIncome, wht float64, allowances []model.Allowance) (*model.TaxCalculationResponse, error)
}
type taxCalculatorService struct {
	taxRepo  repository.TaxRepository
	adminSvc AdminServiceInterface
}

func NewTaxCalculatorService(taxRepo repository.TaxRepository, adminRepo repository.AdminRepository) TaxCalculatorService {
	return &taxCalculatorService{
		taxRepo:  taxRepo,
		adminSvc: NewAdminService(adminRepo),
	}
}

func (s *taxCalculatorService) GetAllCalculations() ([]*model.TaxCalculation, error) {
	return s.taxRepo.GetAllCalculations()
}

func (s *taxCalculatorService) CalculateTax(totalIncome, wht float64, allowances []model.Allowance) (*model.TaxCalculationResponse, error) {
	config, err := s.adminSvc.GetConfig()
	if err != nil {
		return nil, err
	}

	// Set default values if not provided
	personalAllowance := config.PersonalDeduction
	donation := 0.0
	kReceipt := 0.0

	for _, allowance := range allowances {
		switch allowance.AllowanceType {
		case "donation":
			donation = math.Min(allowance.Amount, 100000)
		case "k-receipt":
			kReceipt = math.Min(allowance.Amount, config.KReceipt)
		}
	}

	taxableIncome := totalIncome - personalAllowance - donation - kReceipt

	var tax float64
	switch {
	case taxableIncome <= 150000:
		tax = 0
	case taxableIncome <= 500000:
		tax = (taxableIncome - 150000) * 0.1
	case taxableIncome <= 1000000:
		tax = 35000 + (taxableIncome-500000)*0.15
	case taxableIncome <= 2000000:
		tax = 110000 + (taxableIncome-1000000)*0.2
	default:
		tax = 310000 + (taxableIncome-2000000)*0.35
	}

	taxPayable := math.Max(tax-wht, 0)

	var taxRefund float64
	if tax < wht {
		taxRefund = wht - tax
	}

	taxLevel := []model.TaxRate{
		{Level: "0-150,000", Tax: 0},
		{Level: "150,001-500,000", Tax: 0},
		{Level: "500,001-1,000,000", Tax: 0},
		{Level: "1,000,001-2,000,000", Tax: 0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0},
	}

	switch {
	case taxableIncome <= 150000:
		taxLevel[0].Tax = taxPayable
	case taxableIncome <= 500000:
		taxLevel[1].Tax = taxPayable
	case taxableIncome <= 1000000:
		taxLevel[2].Tax = taxPayable
	case taxableIncome <= 2000000:
		taxLevel[3].Tax = taxPayable
	default:
		taxLevel[4].Tax = taxPayable
	}

	taxCalculation := &model.TaxCalculation{
		TotalIncome:       totalIncome,
		WHT:               wht,
		PersonalAllowance: personalAllowance,
		Donation:          donation,
		KReceipt:          kReceipt,
		Tax:               tax,
		TaxPayable:        taxPayable,
		TaxRefund:         taxRefund,
		TaxLevel:          taxLevel,
	}

	err = s.taxRepo.Save(taxCalculation)
	if err != nil {
		return nil, err
	}

	taxResponse := &model.TaxCalculationResponse{
		TaxLevel: taxLevel,
	}

	if taxPayable > 0 {
		taxResponse.Tax = &taxPayable
	} else if taxRefund > 0 {
		taxResponse.TaxRefund = &taxRefund
	}

	return taxResponse, nil
}
