package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"

	dfErrors "github.com/skuralll/dfeconomy/errors"
)

// /economy

type EconomyCommand struct {
	*BaseCommand
}

func (c EconomyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	o.Printf("help: TODO")
}

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
	}
	// get target uuid
	tn := e.Username.LoadOr(p.Name())
	var uid uuid.UUID
	if tn == p.Name() {
		uid = p.UUID()
	} else {
		// get uuid by name
		var err error
		uid, err = e.svc.GetUUIDByName(context.Background(), tn)
		if err != nil {
			o.Error("Player not found: " + tn)
			return
		}
	}
	// get balance
	amount, err := e.svc.GetBalance(context.Background(), uid)
	if err != nil {
		o.Error("Failed to get balance")
		return
	}
	o.Printf("Balance of %s: %.2f", tn, amount)
}

// /economy set <target> <amount>

type EconomySetCommand struct {
	*BaseCommand
	SubCmd   cmd.SubCommand `cmd:"set" help:"Set the balance of a player."`
	Username string         `cmd:"username"`
	Amount   float64        `cmd:"amount"`
}

func (e EconomySetCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
		return
	}
	// validate amount
	if e.Amount < 0 {
		o.Error("Amount must be at least 0")
		return
	}
	// get target uuid
	var uid uuid.UUID
	if e.Username == p.Name() {
		uid = p.UUID()
	} else {
		var err error
		uid, err = e.svc.GetUUIDByName(context.Background(), e.Username)
		if err != nil {
			o.Error("Player not found: " + e.Username)
			return
		}
	}
	// set balance
	err := e.svc.SetBalance(context.Background(), uid, e.Username, float64(e.Amount))
	if err != nil {
		o.Error("Failed to set balance: " + err.Error())
		return
	}
	o.Printf("Set balance of %s to %.2f", e.Username, e.Amount)
}

// /economy top <page>

const itemCount int = 10 // Number of items per page

type EconomyTopCommand struct {
	*BaseCommand
	SubCmd cmd.SubCommand `cmd:"top" help:"Show the top players by balance."`
	Page   int            `cmd:"page" help:"The page to show."`
}

func (e EconomyTopCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	_, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
		return
	}
	entries, err := e.svc.GetTopBalances(context.Background(), e.Page, itemCount)
	if err != nil {
		if errors.Is(err, dfErrors.ErrValueMustBeAtLeastOne) {
			o.Error("Size must be at least 1")
		} else if errors.Is(err, dfErrors.ErrPageNotFound) {
			o.Error("Page not found")
		} else {
			o.Error("Failed to get top balances by internal error")
		}
		return
	}
	for i, entry := range entries {
		o.Printf("#%d %s: %.2f", (e.Page-1)*itemCount+i+1, entry.Name, entry.Money)
	}
}

// /economy pay <target> <amount>

type EconomyPayCommand struct {
	*BaseCommand
	SubCmd   cmd.SubCommand `cmd:"pay" help:"Pay a player."`
	Username string         `cmd:"username"`
	Amount   float64        `cmd:"amount"`
}

func (e EconomyPayCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
		return
	}
	// get target uuid
	tuid, err := e.svc.GetUUIDByName(context.Background(), e.Username)
	if err != nil {
		o.Error("Player not found: " + e.Username)
		return
	}
	err = e.svc.TransferBalance(context.Background(), p.UUID(), tuid, e.Amount)
	if err != nil {
		switch {
		case errors.Is(err, dfErrors.ErrValueMustBeAtLeastOne):
			o.Error("Amount must be at least 1")
			return
		case errors.Is(err, dfErrors.ErrInsufficientFunds):
			o.Error("Insufficient funds")
			return
		case errors.Is(err, dfErrors.ErrUnknownPlayer):
			o.Error("Target player not found: " + e.Username)
			return
		default:
			o.Error("Failed to pay by internal error")
			slog.Error("Failed to pay", "error", err, "from", p.Name(), "to", e.Username, "amount", e.Amount)
			return
		}
	}
	// success
	o.Printf("You paid %.2f to %s", e.Amount, e.Username)
}

// Validation
var _ cmd.Runnable = (*EconomyCommand)(nil)
var _ cmd.Runnable = (*EconomyBalanceCommand)(nil)
var _ cmd.Runnable = (*EconomySetCommand)(nil)
var _ cmd.Runnable = (*EconomyTopCommand)(nil)
var _ cmd.Runnable = (*EconomyPayCommand)(nil)
