package sqlite_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	_ "modernc.org/sqlite" // SQLite ドライバ

	"github.com/skuralll/dfeconomy/backend/sqlite"
	"github.com/skuralll/dfeconomy/economy"
)

// テスト用サービス生成（インメモリ DB）
func newSvc(t *testing.T) (*sqliteSvcWrapper, context.Context) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	svc := sqlite.New(db)
	ctx := context.Background()

	return &sqliteSvcWrapper{svc}, ctx
}

// ラッパで *testing.T を持たせると呼び出しが簡単
type sqliteSvcWrapper struct{ economy.Economy }

func TestBalance_NotFound(t *testing.T) {
	svc, ctx := newSvc(t)
	_, err := svc.Balance(ctx, uuid.New())
	if !errors.Is(err, economy.ErrUnknownPlayer) {
		t.Fatalf("want ErrUnknownPlayer, got %v", err)
	}
}

func TestSetAndBalance(t *testing.T) {
	svc, ctx := newSvc(t)

	id := uuid.New()
	name := "Alice"
	if err := svc.Set(ctx, id, &name, 123.45); err != nil {
		t.Fatalf("set: %v", err)
	}

	got, err := svc.Balance(ctx, id)
	if err != nil || got != 123.45 {
		t.Fatalf("balance = %v, err = %v", got, err)
	}
}

func TestTopPagination(t *testing.T) {
	svc, ctx := newSvc(t)

	// 3 人分データ投入
	ids := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	amts := []float64{100, 300, 200}

	for i, id := range ids {
		name := fmt.Sprintf("P%d", i)
		_ = svc.Set(ctx, id, &name, amts[i])
	}

	// page=1, size=2 なら 300,200 の順
	list, err := svc.Top(ctx, 1, 2)
	if err != nil {
		t.Fatalf("top: %v", err)
	}

	want := []float64{300, 200}
	for i, e := range list {
		if e.Money != want[i] {
			t.Fatalf("rank %d want %v got %v", i, want[i], e.Money)
		}
	}
}

func TestNegativeAmount(t *testing.T) {
	svc, ctx := newSvc(t)
	err := svc.Set(ctx, uuid.New(), nil, -10)
	if !errors.Is(err, economy.ErrNegativeAmount) {
		t.Fatalf("expected ErrNegativeAmount, got %v", err)
	}
}
