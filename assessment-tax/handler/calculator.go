// Calculator
package handler

import (
	"net/http"

	"github.com/LGROW101/assessment-tax/model"
	"github.com/LGROW101/assessment-tax/service"
	"github.com/labstack/echo/v4"
)

type CalculatorHandler struct {
	taxCalculatorService service.TaxCalculatorService
}

func NewCalculatorHandler(taxCalculatorService service.TaxCalculatorService) *CalculatorHandler {
	return &CalculatorHandler{
		taxCalculatorService: taxCalculatorService,
	}
}

type TaxResponse struct {
	TaxRefund *float64        `json:"taxRefund,omitempty"`
	Tax       *float64        `json:"tax,omitempty"`
	TaxLevel  []model.TaxRate `json:"taxLevel,omitempty"`
}

type CalculateTaxRequest struct {
	TotalIncome     float64           `json:"totalIncome"`
	WHT             float64           `json:"wht"`
	Allowances      []model.Allowance `json:"allowances"`
	IncludeTaxLevel bool              `json:"includeTaxLevel"`
}

func (h *CalculatorHandler) CalculateTax(c echo.Context) error {
	var req CalculateTaxRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.TotalIncome < 0 || req.WHT < 0 || len(req.Allowances) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	taxCalculationResponse, err := h.taxCalculatorService.CalculateTax(req.TotalIncome, req.WHT, req.Allowances)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := TaxResponse{}

	if taxCalculationResponse.TaxRefund != nil && *taxCalculationResponse.TaxRefund > 0 {
		response.TaxRefund = taxCalculationResponse.TaxRefund
	} else if taxCalculationResponse.Tax != nil && *taxCalculationResponse.Tax > 0 {
		response.Tax = taxCalculationResponse.Tax
	}

	if req.IncludeTaxLevel {
		response.TaxLevel = taxCalculationResponse.TaxLevel
	}

	return c.JSON(http.StatusOK, response)
}
func (h *CalculatorHandler) GetAllCalculations(c echo.Context) error {
	taxCalculations, err := h.taxCalculatorService.GetAllCalculations()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, taxCalculations)
}
