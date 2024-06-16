// model/admin.go
package model

import (
	"errors"
	"time"
)

type AdminRequest struct {
	PersonalDeduction *float64 `json:"personalDeduction"`
	KReceipt          *float64 `json:"k_receipt"`
}

type AdminConfig struct {
	ID                uint      `json:"ID,omitempty" gorm:"primaryKey" db:"id"`
	PersonalDeduction float64   `json:"PersonalDeduction,omitempty" db:"personal_deduction"`
	KReceipt          float64   `json:"KReceipt,omitempty" db:"k_receipt"`
	CreatedAt         time.Time `json:"-" db:"created_at"`
	UpdatedAt         time.Time `json:"-" db:"updated_at"`
}

type AdminResponse struct {
	PersonalDeduction float64 `json:"personalDeduction,omitempty"`
	KReceipt          float64 `json:"KReceipt,omitempty"`
}

func (c *AdminConfig) Validate() error {
	if c.PersonalDeduction != 0 {
		if c.PersonalDeduction <= 0 {
			return errors.New("personal deduction must be positive")
		}
	}
	if c.KReceipt != 0 {
		if c.KReceipt <= 0 {
			return errors.New("k-receipt must be positive")
		}
	}
	return nil
}
