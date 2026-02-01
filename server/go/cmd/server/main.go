package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hymaia/lebonpoint-app/server/go/internal/application/usecase"
	"github.com/hymaia/lebonpoint-app/server/go/internal/infrastructure/http/router"
	"github.com/hymaia/lebonpoint-app/server/go/internal/infrastructure/persistence"
	"github.com/hymaia/lebonpoint-app/server/go/internal/shared/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := persistence.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping database to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Printf("Connected to database: %s", cfg.DatabasePath)

	// Initialize repository
	itemRepo := persistence.NewSqliteItemRepository(db.GetDB())

	// Initialize use cases
	createUC := usecase.NewCreateItemUseCase(itemRepo)
	getUC := usecase.NewGetItemUseCase(itemRepo)
	listUC := usecase.NewListItemsUseCase(itemRepo)
	updateUC := usecase.NewUpdateItemUseCase(itemRepo)
	patchUC := usecase.NewPatchItemUseCase(itemRepo)
	deleteUC := usecase.NewDeleteItemUseCase(itemRepo)
	searchUC := usecase.NewSearchItemsUseCase(itemRepo)

	// Setup handlers
	itemHandler := router.SetupHandlers(
		createUC,
		getUC,
		listUC,
		updateUC,
		patchUC,
		deleteUC,
		searchUC,
	)

	// Setup router
	r := router.NewRouter(itemHandler)
	server := &http.Server{
		Addr:    cfg.GetAddr(),
		Handler: r.GetEngine(),
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.GetAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Server exited successfully")
}
