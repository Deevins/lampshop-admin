package repository

import (
	"context"
	"github.com/Deevins/lampshop-admin-backend/internal/repository/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AttributeRepository описывает методы получения категорий и атрибутов.
type AttributeRepository interface {
	ListCategories(ctx context.Context) ([]*sql.Category, error)
	ListAttributesByCategory(ctx context.Context, categoryID uuid.UUID) ([]*sql.AttributeOption, error)
}

// attributeRepoPG — реализация через sqlc/pgx.
type attributeRepoPG struct {
	queries *sql.Queries
}

func NewAttributeRepository(pool *pgxpool.Pool) AttributeRepository {
	return &attributeRepoPG{
		queries: sql.New(pool),
	}
}

func (r *attributeRepoPG) ListCategories(ctx context.Context) ([]*sql.Category, error) {
	return r.queries.ListCategories(ctx)
}

func (r *attributeRepoPG) ListAttributesByCategory(ctx context.Context, categoryID uuid.UUID) ([]*sql.AttributeOption, error) {
	return r.queries.ListAttributeOptionsByCategory(ctx, categoryID)
}
