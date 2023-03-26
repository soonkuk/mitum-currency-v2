package cmds

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ProtoconNet/mitum-currency/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type KeyNewCommand struct {
	baseCommand
	Seed    string `arg:"" name:"seed" optional:"" help:"seed for generating key"`
	KeyType string `help:"select btc or ether" default:"btc"`
}

func NewKeyNewCommand() KeyNewCommand {
	cmd := NewbaseCommand()
	return KeyNewCommand{
		baseCommand: *cmd,
	}
}

func (cmd *KeyNewCommand) Run(pctx context.Context) error {
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	cmd.log.Debug().
		Str("seed", cmd.Seed).
		Msg("flags")

	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	var key base.Privatekey

	switch {
	case len(cmd.Seed) > 0:
		if len(strings.TrimSpace(cmd.Seed)) < 1 {
			cmd.log.Warn().Msg("seed consists with empty spaces")
		}
		if len(cmd.KeyType) > 0 && cmd.KeyType == "ether" {
			i, err := currency.NewMEPrivatekeyFromSeed(cmd.Seed)
			if err != nil {
				return err
			}
			key = i
		} else {
			i, err := base.NewMPrivatekeyFromSeed(cmd.Seed)
			if err != nil {
				return err
			}
			key = i
		}

	default:
		if len(cmd.KeyType) > 0 && cmd.KeyType == "ether" {
			key = currency.NewMEPrivatekey()
		} else {
			key = base.NewMPrivatekey()
		}
	}

	o := struct {
		PrivateKey base.PKKey  `json:"privatekey"` //nolint:tagliatelle //...
		Publickey  base.PKKey  `json:"publickey"`
		Hint       interface{} `json:"hint,omitempty"`
		Seed       string      `json:"seed"`
		Type       string      `json:"type"`
	}{
		Seed:       cmd.Seed,
		PrivateKey: key,
		Publickey:  key.Publickey(),
		Type:       "privatekey",
	}

	if hinter, ok := (interface{})(key).(hint.Hinter); ok {
		o.Hint = hinter.Hint()
	}

	b, err := util.MarshalJSONIndent(o)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(os.Stdout, string(b))

	return nil
}
