package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hymaia/lebonpoint-app/server/go/internal/application/dto"
	"github.com/hymaia/lebonpoint-app/server/go/internal/application/usecase"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/entity"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/repository"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/valueobject"
	apperrors "github.com/hymaia/lebonpoint-app/server/go/internal/shared/errors"
)

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	createUC *usecase.CreateItemUseCase
	getUC    *usecase.GetItemUseCase
	listUC   *usecase.ListItemsUseCase
	updateUC *usecase.UpdateItemUseCase
	patchUC  *usecase.PatchItemUseCase
	deleteUC *usecase.DeleteItemUseCase
	searchUC *usecase.SearchItemsUseCase
}

// NewItemHandler creates a new ItemHandler
func NewItemHandler(
	createUC *usecase.CreateItemUseCase,
	getUC *usecase.GetItemUseCase,
	listUC *usecase.ListItemsUseCase,
	updateUC *usecase.UpdateItemUseCase,
	patchUC *usecase.PatchItemUseCase,
	deleteUC *usecase.DeleteItemUseCase,
	searchUC *usecase.SearchItemsUseCase,
) *ItemHandler {
	return &ItemHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		patchUC:  patchUC,
		deleteUC: deleteUC,
		searchUC: searchUC,
	}
}

// ListItems handles GET /v1/items
func (h *ItemHandler) ListItems(c *gin.Context) {
	params, err := h.parseQueryParams(c)
	if err != nil {
		c.Error(err)
		return
	}

	page, err := h.searchUC.Execute(c.Request.Context(), params)
	if err != nil {
		c.Error(err)
		return
	}

	response := h.toListItemsResponse(page, params.GetLimit())
	c.JSON(http.StatusOK, response)
}

// CreateItem handles POST /v1/items
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationError("Invalid request body", map[string]interface{}{
			"details": err.Error(),
		}))
		return
	}

	item, err := h.createUC.Execute(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, h.toItemResponse(item))
}

// GetItem handles GET /v1/items/:id
func (h *ItemHandler) GetItem(c *gin.Context) {
	id, err := h.parseItemID(c)
	if err != nil {
		c.Error(err)
		return
	}

	item, err := h.getUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, h.toItemResponse(item))
}

// UpdateItem handles PUT /v1/items/:id
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	id, err := h.parseItemID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationError("Invalid request body", map[string]interface{}{
			"details": err.Error(),
		}))
		return
	}

	item, err := h.updateUC.Execute(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, h.toItemResponse(item))
}

// PatchItem handles PATCH /v1/items/:id
func (h *ItemHandler) PatchItem(c *gin.Context) {
	id, err := h.parseItemID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.PatchItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationError("Invalid request body", map[string]interface{}{
			"details": err.Error(),
		}))
		return
	}

	item, err := h.patchUC.Execute(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, h.toItemResponse(item))
}

// DeleteItem handles DELETE /v1/items/:id
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id, err := h.parseItemID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.deleteUC.Execute(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Health handles GET /health
func (h *ItemHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339Nano),
	})
}

// Helper methods

func (h *ItemHandler) parseItemID(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, apperrors.NewBadRequestError("invalid item ID")
	}
	if id <= 0 {
		return 0, apperrors.NewBadRequestError("item ID must be positive")
	}
	return id, nil
}

func (h *ItemHandler) parseQueryParams(c *gin.Context) (*dto.QueryParams, error) {
	params := &dto.QueryParams{}
	details := make(map[string]interface{})

	// Parse filters
	if status := c.Query("status"); status != "" {
		if _, err := valueobject.ParseItemStatus(status); err != nil {
			details["status"] = "status must be one of: draft, active, reserved, sold, archived"
		} else {
			params.Status = &status
		}
	}
	if category := c.Query("category"); category != "" {
		params.Category = &category
	}
	if city := c.Query("city"); city != "" {
		params.City = &city
	}
	if postalCode := queryFirst(c, "postal_code", "postalCode"); postalCode != "" {
		params.PostalCode = &postalCode
	}
	if isFeatured := queryFirst(c, "is_featured", "isFeatured"); isFeatured != "" {
		if val, err := strconv.ParseBool(isFeatured); err == nil {
			params.IsFeatured = &val
		} else {
			details["is_featured"] = "is_featured must be a boolean"
		}
	}
	if deliveryAvailable := queryFirst(c, "delivery_available", "deliveryAvailable"); deliveryAvailable != "" {
		if val, err := strconv.ParseBool(deliveryAvailable); err == nil {
			params.DeliveryAvailable = &val
		} else {
			details["delivery_available"] = "delivery_available must be a boolean"
		}
	}

	// Parse sort and pagination
	if sort := c.Query("sort"); sort != "" {
		params.Sort = &sort
	}
	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			if val < 1 || val > 100 {
				details["limit"] = "limit must be between 1 and 100"
			} else {
				params.Limit = &val
			}
		} else {
			details["limit"] = "limit must be an integer"
		}
	}
	if cursor := c.Query("cursor"); cursor != "" {
		params.Cursor = &cursor
	}

	// Parse search query
	if query := c.Query("q"); query != "" {
		params.Query = &query
	}

	if len(details) > 0 {
		return nil, apperrors.NewValidationError("Validation failed", details)
	}

	return params, nil
}

func queryFirst(c *gin.Context, keys ...string) string {
	for _, key := range keys {
		if value := c.Query(key); value != "" {
			return value
		}
	}
	return ""
}

func (h *ItemHandler) toItemResponse(item *entity.Item) *dto.ItemResponse {
	return &dto.ItemResponse{
		ID:                item.ID,
		Title:             item.Title,
		Description:       item.Description,
		PriceCents:        item.PriceCents,
		Category:          item.Category,
		Condition:         item.Condition,
		Status:            item.Status,
		IsFeatured:        item.IsFeatured,
		City:              item.City,
		PostalCode:        item.PostalCode,
		Country:           item.Country,
		DeliveryAvailable: item.DeliveryAvailable,
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
		PublishedAt:       item.PublishedAt,
		Images:            item.Images,
	}
}

func (h *ItemHandler) toListItemsResponse(page *repository.ItemPage, limit int) *dto.ListItemsResponse {
	items := make([]*dto.ItemResponse, len(page.Items))
	for i, item := range page.Items {
		items[i] = h.toItemResponse(item)
	}

	return &dto.ListItemsResponse{
		Items:      items,
		NextCursor: page.NextCursor,
		TotalCount: page.TotalCount,
		Limit:      limit,
	}
}
