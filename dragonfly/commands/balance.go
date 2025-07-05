package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

// /economy balance <target>

type EconomyBalanceCommand struct {
	*BaseCommand
	SubCmd   cmd.SubCommand       `cmd:"balance" help:"Displays the balance of a player."`
	Username cmd.Optional[string] `cmd:"username"`
}

func (e EconomyBalanceCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
		return
	}

	// 即座にフィードバック
	o.Printf("Fetching balance...")

	go func() {
		// タイムアウト付きでDB処理
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// get target uuid
		tn := e.Username.LoadOr(p.Name())
		var uid uuid.UUID
		if tn == p.Name() {
			uid = p.UUID()
		} else {
			// get uuid by name
			var err error
			uid, err = e.svc.GetUUIDByName(ctx, tn)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					p.Message("§c[Error] Request timeout")
				} else {
					p.Message("§c[Error] Player not found: " + tn)
				}
				return
			}
		}

		// get balance
		amount, err := e.svc.GetBalance(ctx, uid)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				p.Message("§c[Error] Request timeout")
			} else {
				p.Message("§c[Error] Failed to get balance")
			}
			return
		}
		// 結果を直接送信
		p.Message(fmt.Sprintf("§a[Balance] %s: %.2f", tn, amount))
	}()
}

// Validation
var _ cmd.Runnable = (*EconomyBalanceCommand)(nil)
