package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

	"sweng-task/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Errors
var (
	ErrLineItemNotFound = errors.New("line item not found")
)

// LineItemService provides operations for line items
type LineItemService struct {
	db_pool *pgxpool.Pool
	log     *zap.SugaredLogger
}

// NewLineItemService creates a new LineItemService
func NewLineItemService(log *zap.SugaredLogger, pool *pgxpool.Pool) *LineItemService {
	return &LineItemService{
		db_pool: pool,
		log:     log,
	}
}

// Create creates a new line item
func (s *LineItemService) Create(item model.LineItemCreate) (*model.LineItem, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	now := time.Now()
	lineItemID := "li_" + uuid.New().String()

	lineItem := &model.LineItem{
		ID:           lineItemID,
		Name:         item.Name,
		AdvertiserID: item.AdvertiserID,
		Bid:          item.Bid,
		Budget:       item.Budget,
		Placement:    item.Placement,
		Categories:   item.Categories,
		Keywords:     item.Keywords,
		Status:       model.LineItemStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	statement := `INSERT INTO line_items 
    			(id, name, advertiser_id, bid, budget, placement, categories, keywords, status, created_at, updated_at) 
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`

	res, err := s.db_pool.Exec(ctx, statement, lineItemID,
		item.Name, item.AdvertiserID, item.Bid, item.Budget, item.Placement, item.Categories, item.Keywords,
		model.LineItemStatusActive, now, now)

	if err != nil {
		fmt.Println("Error inserting line item:", err)
		s.log.Errorw("Failed to create line item",
			"error", err,
			"line_item", lineItem,
			"postgres result", res,
		)
		return nil, err
	}

	s.log.Infow("Line item created",
		"id", lineItemID,
		"name", item.Name,
		"advertiser_id", item.AdvertiserID,
		"placement", item.Placement,
		"postgres result", res,
	)

	return lineItem, nil
}

// GetByID retrieves a line item by ID
func (s *LineItemService) GetByID(id string) (*model.LineItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	row, err := s.db_pool.Query(ctx,
		`SELECT * FROM line_items WHERE id=$1;`,
		id)
	defer row.Close()

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrLineItemNotFound
		}
		s.log.Errorw("Failed to retrieve line item",
			"error", err,
			"id", id,
		)
		return nil, err
	}

	isNextRow := row.Next()

	if !isNextRow {
		return nil, ErrLineItemNotFound
	}

	// item, err := pgx.RowTo[model.LineItem](row)

	var item model.LineItem
	err = row.Scan(
		&item.ID, &item.Name, &item.AdvertiserID, &item.Bid, &item.Budget, &item.Placement, &item.Categories,
		&item.Keywords, &item.Status, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		s.log.Errorw("Failed to parse line item",
			"error", err,
			"id", id,
		)
		return nil, err
	}

	return &item, err
}

// GetAll retrieves all line items, optionally filtered by advertiser ID and placement
func (s *LineItemService) GetAll(advertiserID, placement, category, keyword, status string) ([]*model.LineItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	var result []*model.LineItem

	statement := `
	SELECT * 
	FROM line_items 
    WHERE true
    	AND (($1='') OR advertiser_id=$1) 
    	AND (($2='') OR placement=$2)
		AND (($3='') OR category=$3)
		AND (($4='') OR keyword=$4)
		AND (($5='') OR status=$5);`

	rows, err := s.db_pool.Query(ctx, statement, advertiserID, placement, category, keyword, status)
	defer rows.Close()

	if err != nil {
		fmt.Println("Error retrieving line items:", err)
		if err == pgx.ErrNoRows {
			return nil, ErrLineItemNotFound
		}
		s.log.Errorw("Failed to retrieve line items",
			"error", err,
			"advertiser_id", advertiserID,
			"placement", placement,
		)
		return nil, err
	}

	for rows.Next() {
		var item model.LineItem
		err := rows.Scan(&item.ID, &item.Name, &item.AdvertiserID, &item.Bid, &item.Budget, &item.Placement, &item.Categories, &item.Keywords, &item.Status, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			fmt.Println("Error scanning line item:", err)
			s.log.Errorw("Failed to parse line item",
				"error", err,
				"advertiser_id", advertiserID,
				"placement", placement,
			)
			return nil, err
		}
		result = append(result, &item)
	}

	return result, nil
}

// FindMatchingLineItems finds line items matching the given placement and filters
// This method will be used by the AdService when implementing the ad selection logic
func (s *LineItemService) FindMatchingLineItems(placement string, category, keyword string) ([]*model.LineItem, error) {
	_, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	var result []*model.LineItem

	/*	for _, item := range s.items {
			// Skip items not matching the placement or not active
			if item.Placement != placement || item.Status != model.LineItemStatusActive {
				continue
			}

			// Apply category filter if specified
			if category != "" {
				categoryFound := false
				for _, cat := range item.Categories {
					if cat == category {
						categoryFound = true
						break
					}
				}
				if !categoryFound {
					continue
				}
			}

			// Apply keyword filter if specified
			if keyword != "" {
				keywordFound := false
				for _, kw := range item.Keywords {
					if kw == keyword {
						keywordFound = true
						break
					}
				}
				if !keywordFound {
					continue
				}
			}

			result = append(result, item)
		}
	*/
	return result, nil
}
