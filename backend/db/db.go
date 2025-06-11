package db

import (
	"context"

	"github.com/google/uuid"
)

type EconomyEntry struct {
	UUID  uuid.UUID // Player’s UUID
	Name  string    // Display name
	Money float64   // Balance
}

type DB interface {
	// 残高取得
	Balance(ctx context.Context, id uuid.UUID) (float64, error)
	// 残高設定
	Set(ctx context.Context, id uuid.UUID, name *string, amount float64) error
	// 残高ランキング取得
	Top(ctx context.Context, page, size int) ([]EconomyEntry, error)
}
