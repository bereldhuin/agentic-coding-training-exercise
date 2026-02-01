package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/entity"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/repository"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/valueobject"
	"github.com/hymaia/lebonpoint-app/server/go/internal/shared/cursor"
	apperrors "github.com/hymaia/lebonpoint-app/server/go/internal/shared/errors"
	"github.com/jmoiron/sqlx"
)

const (
	itemsTable = "items"
)

// sqliteItemRepository implements ItemRepository using SQLite
type sqliteItemRepository struct {
	db *sqlx.DB
}

// NewSqliteItemRepository creates a new SQLite item repository
func NewSqliteItemRepository(db *sqlx.DB) repository.ItemRepository {
	return &sqliteItemRepository{db: db}
}

// Create creates a new item
func (r *sqliteItemRepository) Create(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	imagesJSON, err := json.Marshal(item.Images)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Sprintf("failed to marshal images: %v", err))
	}

	query := `
		INSERT INTO items (
			title, description, price_cents, category, condition, status,
			is_featured, city, postal_code, country, delivery_available,
			created_at, updated_at, published_at, images
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	// Handle nullable published_at
	var publishedAt *string
	if item.PublishedAt != nil {
		publishedAtStr := item.PublishedAt.Format(time.RFC3339Nano)
		publishedAt = &publishedAtStr
	}

	result, err := r.db.ExecContext(ctx, query,
		item.Title,
		item.Description,
		item.PriceCents,
		item.Category,
		item.Condition.String(),
		item.Status.String(),
		boolToInt(item.IsFeatured),
		item.City,
		item.PostalCode,
		item.Country,
		boolToInt(item.DeliveryAvailable),
		item.CreatedAt.Format(time.RFC3339Nano),
		item.UpdatedAt.Format(time.RFC3339Nano),
		publishedAt,
		string(imagesJSON),
	)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Sprintf("failed to create item: %v", err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Sprintf("failed to get last insert id: %v", err))
	}

	item.ID = id
	return item, nil
}

// FindByID retrieves an item by its ID
func (r *sqliteItemRepository) FindByID(ctx context.Context, id int64) (*entity.Item, error) {
	query := `SELECT * FROM items WHERE id = ?`
	var dbItem dbItem
	if err := r.db.GetContext(ctx, &dbItem, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewNotFoundError("Item")
		}
		return nil, apperrors.NewInternalError(fmt.Sprintf("failed to find item: %v", err))
	}

	return scanItem(&dbItem)
}

// FindAll retrieves items with filtering, sorting, and pagination
func (r *sqliteItemRepository) FindAll(ctx context.Context, filters *repository.ItemFilters, sort *repository.SortOptions, limit int, cursorStr string) (*repository.ItemPage, error) {
	// Set defaults
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if sort == nil {
		sort = repository.DefaultSortOptions()
	}

	// Build WHERE clause and args
	whereClause, args := r.buildWhereClause(filters)

	// Parse cursor for pagination
	cursorID := int64(0)
	if cursorStr != "" {
		id, err := cursor.Decode(cursorStr)
		if err != nil {
			return nil, apperrors.NewBadRequestError(fmt.Sprintf("invalid cursor: %v", err))
		}
		cursorID = id
		if whereClause != "" {
			whereClause += " AND id > ?"
		} else {
			whereClause = "WHERE id > ?"
		}
		args = append(args, cursorID)
	}

	// Build ORDER BY clause
	orderBy := r.buildOrderByClause(sort)

	// Build query
	query := fmt.Sprintf(`
		SELECT * FROM items
		%s
		%s
		LIMIT ?
	`, whereClause, orderBy)

	args = append(args, limit+1) // Fetch one extra to check if there are more results

	// Execute query
	var dbItems []dbItem
	if err := r.db.SelectContext(ctx, &dbItems, query, args...); err != nil {
		return nil, apperrors.NewInternalError(fmt.Sprintf("failed to list items: %v", err))
	}

	// Convert to entities
	items := make([]*entity.Item, 0, len(dbItems))
	for _, dbItem := range dbItems {
		item, err := scanItem(&dbItem)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	// Handle pagination
	var nextCursor *string
	if len(items) > limit {
		items = items[:limit] // Remove the extra item
		lastID := items[len(items)-1].ID
		next := cursor.Encode(lastID)
		nextCursor = &next
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM items %s", whereClause)
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, countQuery, args[:len(args)-1]...); err != nil {
		totalCount = len(items) // Fallback if count fails
	}

	return repository.NewItemPage(items, nextCursor, totalCount), nil
}

// Update updates an existing item
func (r *sqliteItemRepository) Update(ctx context.Context, item *entity.Item) error {
	// Check if item exists
	existing, err := r.FindByID(ctx, item.ID)
	if err != nil {
		return err
	}

	// Update timestamp
	item.UpdatedAt = time.Now()
	if existing.Status != item.Status && item.Status == valueobject.ItemStatusActive {
		now := time.Now()
		item.PublishedAt = &now
	} else if existing.Status == valueobject.ItemStatusActive && item.Status != valueobject.ItemStatusActive {
		item.PublishedAt = nil
	}

	imagesJSON, err := json.Marshal(item.Images)
	if err != nil {
		return apperrors.NewInternalError(fmt.Sprintf("failed to marshal images: %v", err))
	}

	query := `
		UPDATE items SET
			title = ?, description = ?, price_cents = ?, category = ?,
			condition = ?, status = ?, is_featured = ?, city = ?,
			postal_code = ?, country = ?, delivery_available = ?,
			updated_at = ?, published_at = ?, images = ?
		WHERE id = ?
	`

	var publishedAt *string
	if item.PublishedAt != nil {
		publishedAtStr := item.PublishedAt.Format(time.RFC3339Nano)
		publishedAt = &publishedAtStr
	}

	result, err := r.db.ExecContext(ctx, query,
		item.Title,
		item.Description,
		item.PriceCents,
		item.Category,
		item.Condition.String(),
		item.Status.String(),
		boolToInt(item.IsFeatured),
		item.City,
		item.PostalCode,
		item.Country,
		boolToInt(item.DeliveryAvailable),
		item.UpdatedAt.Format(time.RFC3339Nano),
		publishedAt,
		string(imagesJSON),
		item.ID,
	)
	if err != nil {
		return apperrors.NewInternalError(fmt.Sprintf("failed to update item: %v", err))
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewInternalError(fmt.Sprintf("failed to get rows affected: %v", err))
	}
	if rows == 0 {
		return apperrors.NewNotFoundError("Item")
	}

	return nil
}

// Delete deletes an item
func (r *sqliteItemRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM items WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return apperrors.NewInternalError(fmt.Sprintf("failed to delete item: %v", err))
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewInternalError(fmt.Sprintf("failed to get rows affected: %v", err))
	}
	if rows == 0 {
		return apperrors.NewNotFoundError("Item")
	}

	return nil
}

// Search performs full-text search
func (r *sqliteItemRepository) Search(ctx context.Context, queryStr string, filters *repository.ItemFilters, sort *repository.SortOptions, limit int, cursorStr string) (*repository.ItemPage, error) {
	if queryStr == "" {
		return r.FindAll(ctx, filters, sort, limit, cursorStr)
	}

	// Set defaults
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if sort == nil {
		sort = repository.DefaultSortOptions()
	}

	// Build FTS5 query
	ftsQuery := fmt.Sprintf("SELECT rowid FROM items_fts WHERE items_fts MATCH ?")

	// Parse cursor
	cursorID := int64(0)
	if cursorStr != "" {
		id, err := cursor.Decode(cursorStr)
		if err != nil {
			return nil, apperrors.NewBadRequestError(fmt.Sprintf("invalid cursor: %v", err))
		}
		cursorID = id
	}

	// Build filter conditions
	filterConditions, filterArgs := r.buildFilterConditions(filters)

	// Build main query
	mainQuery := `
		SELECT i.* FROM items i
		INNER JOIN (%s) fts ON i.id = fts.rowid
		%s
		%s
		LIMIT ?
	`

	whereClause := ""
	args := []interface{}{queryStr}

	if filterConditions != "" {
		whereClause = "WHERE " + filterConditions
		args = append(args, filterArgs...)
	}

	if cursorID > 0 {
		if whereClause != "" {
			whereClause += " AND i.id > ?"
		} else {
			whereClause = "WHERE i.id > ?"
		}
		args = append(args, cursorID)
	}

	// FTS5 uses BM25 ranking by default when using MATCH
	orderBy := "ORDER BY fts.rowid"
	if sort != nil {
		orderBy = r.buildOrderByClause(sort)
		// For FTS, we might want to rank by relevance first
		orderBy = "ORDER BY rank, i.id"
	}

	query := fmt.Sprintf(mainQuery, ftsQuery, whereClause, orderBy)
	args = append(args, limit+1)

	// Execute query
	var dbItems []dbItem
	if err := r.db.SelectContext(ctx, &dbItems, query, args...); err != nil {
		return nil, apperrors.NewInternalError(fmt.Sprintf("failed to search items: %v", err))
	}

	// Convert to entities
	items := make([]*entity.Item, 0, len(dbItems))
	for _, dbItem := range dbItems {
		item, err := scanItem(&dbItem)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	// Handle pagination
	var nextCursor *string
	if len(items) > limit {
		items = items[:limit]
		lastID := items[len(items)-1].ID
		next := cursor.Encode(lastID)
		nextCursor = &next
	}

	return repository.NewItemPage(items, nextCursor, len(items)), nil
}

// Helper functions

type dbItem struct {
	ID                int64  `db:"id"`
	Title             string `db:"title"`
	Description       string `db:"description"`
	PriceCents        int    `db:"price_cents"`
	Category          string `db:"category"`
	Condition         string `db:"condition"`
	Status            string `db:"status"`
	IsFeatured        int    `db:"is_featured"`
	City              string `db:"city"`
	PostalCode        string `db:"postal_code"`
	Country           string `db:"country"`
	DeliveryAvailable int    `db:"delivery_available"`
	CreatedAt         string `db:"created_at"`
	UpdatedAt         string `db:"updated_at"`
	PublishedAt       string `db:"published_at"`
	Images            string `db:"images"`
}

func scanItem(dbItem *dbItem) (*entity.Item, error) {
	condition, err := valueobject.ParseItemCondition(dbItem.Condition)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Sprintf("invalid condition in database: %s", dbItem.Condition))
	}

	status, err := valueobject.ParseItemStatus(dbItem.Status)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Sprintf("invalid status in database: %s", dbItem.Status))
	}

	createdAt, err := time.Parse(time.RFC3339Nano, dbItem.CreatedAt)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Sprintf("invalid created_at in database: %s", dbItem.CreatedAt))
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, dbItem.UpdatedAt)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Sprintf("invalid updated_at in database: %s", dbItem.UpdatedAt))
	}

	var publishedAt *time.Time
	if dbItem.PublishedAt != "" {
		parsed, err := time.Parse(time.RFC3339Nano, dbItem.PublishedAt)
		if err != nil {
			return nil, apperrors.NewInternalError(fmt.Sprintf("invalid published_at in database: %s", dbItem.PublishedAt))
		}
		publishedAt = &parsed
	}

	var images []valueobject.ItemImage
	if dbItem.Images != "" && dbItem.Images != "[]" {
		if err := json.Unmarshal([]byte(dbItem.Images), &images); err != nil {
			return nil, apperrors.NewInternalError(fmt.Sprintf("failed to parse images: %v", err))
		}
	}
	if images == nil {
		images = []valueobject.ItemImage{}
	}

	return &entity.Item{
		ID:                dbItem.ID,
		Title:             dbItem.Title,
		Description:       dbItem.Description,
		PriceCents:        dbItem.PriceCents,
		Category:          dbItem.Category,
		Condition:         condition,
		Status:            status,
		IsFeatured:        intToBool(dbItem.IsFeatured),
		City:              dbItem.City,
		PostalCode:        dbItem.PostalCode,
		Country:           dbItem.Country,
		DeliveryAvailable: intToBool(dbItem.DeliveryAvailable),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		PublishedAt:       publishedAt,
		Images:            images,
	}, nil
}

func (r *sqliteItemRepository) buildWhereClause(filters *repository.ItemFilters) (string, []interface{}) {
	if filters == nil {
		return "", nil
	}

	conditions, args := r.buildFilterConditions(filters)
	if conditions == "" {
		return "", nil
	}

	return "WHERE " + conditions, args
}

func (r *sqliteItemRepository) buildFilterConditions(filters *repository.ItemFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filters.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, filters.Status.String())
	}

	if filters.Category != nil && *filters.Category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, *filters.Category)
	}

	if filters.MinPriceCents != nil {
		conditions = append(conditions, "price_cents >= ?")
		args = append(args, *filters.MinPriceCents)
	}

	if filters.MaxPriceCents != nil {
		conditions = append(conditions, "price_cents <= ?")
		args = append(args, *filters.MaxPriceCents)
	}

	if filters.City != nil && *filters.City != "" {
		conditions = append(conditions, "city = ?")
		args = append(args, *filters.City)
	}

	if filters.PostalCode != nil && *filters.PostalCode != "" {
		conditions = append(conditions, "postal_code = ?")
		args = append(args, *filters.PostalCode)
	}

	if filters.IsFeatured != nil {
		conditions = append(conditions, "is_featured = ?")
		args = append(args, boolToInt(*filters.IsFeatured))
	}

	if filters.DeliveryAvailable != nil {
		conditions = append(conditions, "delivery_available = ?")
		args = append(args, boolToInt(*filters.DeliveryAvailable))
	}

	return strings.Join(conditions, " AND "), args
}

func (r *sqliteItemRepository) buildOrderByClause(sort *repository.SortOptions) string {
	direction := strings.ToUpper(string(sort.Direction))
	if direction != "ASC" && direction != "DESC" {
		direction = "DESC"
	}

	field := string(sort.Field)
	// Map Go field names to database column names
	columnMap := map[string]string{
		"id":          "id",
		"title":       "title",
		"priceCents":  "price_cents",
		"createdAt":   "created_at",
		"updatedAt":   "updated_at",
		"publishedAt": "published_at",
	}

	column, ok := columnMap[field]
	if !ok {
		column = "created_at" // Default
	}

	return fmt.Sprintf("ORDER BY %s %s", column, direction)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i int) bool {
	return i == 1
}
