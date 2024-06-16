package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/LGROW101/assessment-tax/config"
	"github.com/LGROW101/assessment-tax/handler"
	"github.com/LGROW101/assessment-tax/repository"
	"github.com/LGROW101/assessment-tax/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create repository instances
	taxRepo := repository.NewTaxRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// Create service instances
	taxCalculatorService := service.NewTaxCalculatorService(taxRepo, adminRepo)
	taxCSVService := service.NewTaxCSVService(taxRepo, adminRepo)
	// Create handler instances
	calculatorHandler := handler.NewCalculatorHandler(taxCalculatorService)
	csvHandler := handler.NewCSVHandler(taxCSVService)
	adminHandler := handler.NewAdminHandler(adminRepo)

	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", func(c echo.Context) error {
		// Define your response data
		responseData := map[string]interface{}{
			"message": "Hello, world!",
		}
		return c.JSON(http.StatusOK, responseData)
	})

	e.POST("tax/calculations", calculatorHandler.CalculateTax)

	e.GET("tax/calculations", calculatorHandler.GetAllCalculations)

	e.GET("/admin/deductions", adminHandler.GetConfig)

	e.POST("/admin/deductions", adminHandler.UpdateConfig, middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == cfg.AdminUsername && password == cfg.AdminPassword, nil
	}))

	e.POST("tax/calculations/upload-csv", csvHandler.UploadCSV)

	// Start server with graceful shutdown
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
