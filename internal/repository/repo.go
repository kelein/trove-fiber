package repository

import (
	"context"

	"gorm.io/gorm"
)

const ctxTxKey = "TxKey"

// Transaction standard interface for transaction management
type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// NewTransaction builds a new Transaction instance
func NewTransaction(r *Repository) Transaction { return r }

// Repository layer for database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository instance
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// DB builds a gorm.DB instance with the context
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}

// Transaction allows executing a function within a transaction context
func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}
