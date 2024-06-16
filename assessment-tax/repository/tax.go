// tax
package repository

import (
	"database/sql"

	"github.com/LGROW101/assessment-tax/model"
)

type TaxRepository interface {
	Save(tax *model.TaxCalculation) error
	GetAllCalculations() ([]*model.TaxCalculation, error)
}

type taxRepository struct {
	db *sql.DB
}

func NewTaxRepository(db *sql.DB) TaxRepository {
	return &taxRepository{db: db}
}

func (r *taxRepository) Save(tax *model.TaxCalculation) error {

	query := `
	INSERT INTO tax_calculations (
		totalIncome,
		wht,
		personal_allowance,
		donation,
		k_receipt,
		tax
	) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		query,
		tax.TotalIncome,
		tax.WHT,
		tax.PersonalAllowance,
		tax.Donation,
		tax.KReceipt,
		tax.Tax,
	)

	return err
}

func (r *taxRepository) GetAllCalculations() ([]*model.TaxCalculation, error) {

	var taxCalculations []*model.TaxCalculation

	query := `
		SELECT
			id,
			totalIncome,
			wht,
			personal_allowance,
			donation,
			k_receipt,
			tax,
			created_at
		FROM
			tax_calculations
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var taxCalculation model.TaxCalculation

		err := rows.Scan(
			&taxCalculation.ID,
			&taxCalculation.TotalIncome,
			&taxCalculation.WHT,
			&taxCalculation.PersonalAllowance,
			&taxCalculation.Donation,
			&taxCalculation.KReceipt,
			&taxCalculation.Tax,
			&taxCalculation.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		taxCalculations = append(taxCalculations, &taxCalculation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return taxCalculations, nil
}
