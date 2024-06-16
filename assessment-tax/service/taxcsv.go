package service

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/LGROW101/assessment-tax/repository"
)

type TaxCSVService interface {
	ImportCSV(reader io.Reader) ([]map[string]float64, error)
	CalculateTax(totalIncome, wht, donation float64) (float64, float64, error)
}

type taxCSVService struct {
	taxRepo  repository.TaxRepository
	adminSvc AdminServiceInterface
}

// NewTaxCSVService returns a new instance of TaxCSVService
func NewTaxCSVService(taxRepo repository.TaxRepository, adminRepo repository.AdminRepository) TaxCSVService {
	return &taxCSVService{
		taxRepo:  taxRepo,
		adminSvc: NewAdminService(adminRepo),
	}
}

func (s *taxCSVService) ImportCSV(reader io.Reader) ([]map[string]float64, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields per record

	csvReader.Read()

	var taxes []map[string]float64
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		totalIncome, wht, donation, err := ParseFields(line)
		if err != nil {
			return nil, err
		}

		taxPayable, taxRefund, err := s.CalculateTax(totalIncome, wht, donation)
		if err != nil {
			return nil, err
		}

		taxResult := map[string]float64{
			"totalIncome": totalIncome,
		}

		if taxRefund > 0 {
			taxResult["taxRefund"] = taxRefund
		} else {
			taxResult["tax"] = taxPayable
		}

		taxes = append(taxes, taxResult)
	}

	return taxes, nil
}

func (s *taxCSVService) CalculateTax(totalIncome, wht, donation float64) (float64, float64, error) {

	config, err := s.adminSvc.GetConfig()
	if err != nil {
		return 0, 0, err
	}

	personalAllowance := config.PersonalDeduction

	taxableIncome := totalIncome - personalAllowance - donation

	var tax float64
	switch {
	case taxableIncome <= 0:
		tax = 0
	case taxableIncome <= 150000:
		tax = taxableIncome * 0.05
	case taxableIncome <= 300000:
		tax = 7500 + (taxableIncome-150000)*0.05
	case taxableIncome <= 500000:
		tax = 15000 + (taxableIncome-300000)*0.10
	case taxableIncome <= 750000:
		tax = 35000 + (taxableIncome-500000)*0.15
	case taxableIncome <= 1000000:
		tax = 57500 + (taxableIncome-750000)*0.20
	case taxableIncome <= 2000000:
		tax = 107500 + (taxableIncome-1000000)*0.25
	case taxableIncome <= 5000000:
		tax = 357500 + (taxableIncome-2000000)*0.30
	default:
		tax = 1257500 + (taxableIncome-5000000)*0.35
	}

	// Calculate tax payable
	taxPayable := tax - wht
	var taxRefund float64
	if taxPayable < 0 {
		taxRefund = -taxPayable
		taxPayable = 0
	}

	return taxPayable, taxRefund, nil
}

func ParseFields(fields []string) (float64, float64, float64, error) {
	var totalIncome, wht, donation float64
	var err error

	if len(fields) > 0 {
		totalIncome, err = strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	if len(fields) > 1 {
		wht, err = strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	if len(fields) > 2 {
		donation, err = ParseDonation(strings.TrimSpace(fields[2]))
		if err != nil {
			return 0, 0, 0, err
		}
	}

	return totalIncome, wht, donation, nil
}

func ParseDonation(donationStr string) (float64, error) {
	if strings.HasSuffix(donationStr, "%") {
		percentage, err := strconv.ParseFloat(strings.TrimSuffix(donationStr, "%"), 64)
		if err != nil {
			return 0, err
		}
		return percentage / 100, nil
	}
	donation, err := strconv.ParseFloat(donationStr, 64)
	if err != nil {
		return 0, err
	}
	return donation, nil
}
