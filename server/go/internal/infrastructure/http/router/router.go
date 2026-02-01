package router

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/hymaia/lebonpoint-app/server/go/internal/application/usecase"
	"github.com/hymaia/lebonpoint-app/server/go/internal/infrastructure/http/handler"
	"github.com/hymaia/lebonpoint-app/server/go/internal/infrastructure/http/middleware"
)

// Router wraps the Gin engine
type Router struct {
	engine *gin.Engine
}

// NewRouter creates a new router with all routes configured
func NewRouter(
	itemHandler *handler.ItemHandler,
) *Router {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin engine
	engine := gin.New()

	// Add middleware
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.LoggingMiddleware())
	engine.Use(middleware.ErrorHandler())

	indexPath := resolveClientIndexPath()

	// Root page
	engine.GET("/", func(c *gin.Context) {
		if indexPath == "" {
			c.Status(500)
			return
		}
		if _, err := os.Stat(indexPath); err != nil {
			c.Status(404)
			return
		}
		c.File(indexPath)
	})

	// Health check route
	engine.GET("/health", itemHandler.Health)

	// API v1 routes
	v1 := engine.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.GET("", itemHandler.ListItems)
			items.POST("", itemHandler.CreateItem)
			items.GET("/:id", itemHandler.GetItem)
			items.PUT("/:id", itemHandler.UpdateItem)
			items.PATCH("/:id", itemHandler.PatchItem)
			items.DELETE("/:id", itemHandler.DeleteItem)
		}
	}

	return &Router{
		engine: engine,
	}
}

func resolveClientIndexPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	dir := filepath.Dir(filename)
	for i := 0; i < 10; i++ {
		candidate := filepath.Join(dir, "client", "index.html")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// GetEngine returns the Gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// SetupHandlers creates and wires all handlers
func SetupHandlers(
	createUC *usecase.CreateItemUseCase,
	getUC *usecase.GetItemUseCase,
	listUC *usecase.ListItemsUseCase,
	updateUC *usecase.UpdateItemUseCase,
	patchUC *usecase.PatchItemUseCase,
	deleteUC *usecase.DeleteItemUseCase,
	searchUC *usecase.SearchItemsUseCase,
) *handler.ItemHandler {
	return handler.NewItemHandler(
		createUC,
		getUC,
		listUC,
		updateUC,
		patchUC,
		deleteUC,
		searchUC,
	)
}
