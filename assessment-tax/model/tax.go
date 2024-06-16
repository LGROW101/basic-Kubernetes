package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TaxCalculation struct {
	ID                uint        `gorm:"primaryKey"`
	TotalIncome       float64     `db:"totalIncome"`
	WHT               float64     `db:"wht"`
	PersonalAllowance float64     `db:"personal_allowance"`
	Donation          float64     `db:"donation"`
	KReceipt          float64     `db:"k_receipt"`
	Tax               float64     `db:"tax"`
	TaxPayable        float64     `db:"tax_payable"`
	TaxRefund         float64     `json:"taxRefund"`
	TaxLevel          []TaxRate   `gorm:"-" json:"taxLevel"`
	Allowances        []Allowance `json:"allowances"`
	CreatedAt         time.Time   `json:"createdAt"`
}
type TaxRate struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxCalculationResponse struct {
	Tax       *float64  `json:"tax,omitempty"`
	TaxRefund *float64  `json:"taxRefund,omitempty"`
	TaxLevel  []TaxRate `json:"taxLevel,omitempty"`
}

func (t *TaxCalculation) BeforeSave(tx *gorm.DB) (err error) {

	if t.TotalIncome < 0 {
		return errors.New("income must be positive")
	}

	return nil
}
