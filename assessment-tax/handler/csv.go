package handler

import (
	"net/http"

	"github.com/LGROW101/assessment-tax/service"
	"github.com/labstack/echo/v4"
)

type CSVHandler struct {
	taxCSVService service.TaxCSVService
}

func NewCSVHandler(taxCSVService service.TaxCSVService) *CSVHandler {
	return &CSVHandler{
		taxCSVService: taxCSVService,
	}
}

func (h *CSVHandler) UploadCSV(c echo.Context) error {
	file, err := c.FormFile("taxFile")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	taxes, err := h.taxCSVService.ImportCSV(src)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"taxes": taxes})
}
