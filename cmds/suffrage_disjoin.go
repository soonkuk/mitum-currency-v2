package cmds

import (
	"context"

	isaacoperation "github.com/ProtoconNet/mitum-currency/isaac"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

type SuffrageDisjoinCommand struct {
	baseCommand
	OperationFlags
	Node  AddressFlag `arg:"" name:"node" help:"node address" required:"true"`
	Start base.Height `arg:"" name:"height" help:"block height" required:"true"`
	node  base.Address
}

func NewSuffrageDisjoinCommand() SuffrageDisjoinCommand {
	cmd := NewbaseCommand()
	return SuffrageDisjoinCommand{
		baseCommand: *cmd,
	}
}

func (cmd *SuffrageDisjoinCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.encs
	enc = cmd.enc

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	var op base.Operation
	if i, err := cmd.createOperation(); err != nil {
		return errors.Wrap(err, "failed to create suffrage-disjoin operation")
	} else if err := i.IsValid([]byte(cmd.OperationFlags.NetworkID)); err != nil {
		return errors.Wrap(err, "invalid suffrage-disjoin operation")
	} else {
		cmd.log.Debug().Interface("operation", i).Msg("operation loaded")

		op = i
	}

	PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *SuffrageDisjoinCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Node.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid node format, %q", cmd.Node.String())
	}
	cmd.node = a

	return nil
}

func (cmd *SuffrageDisjoinCommand) createOperation() (isaacoperation.SuffrageDisjoin, error) {
	fact := isaacoperation.NewSuffrageDisjoinFact([]byte(cmd.Token), cmd.node, cmd.Start)

	op := isaacoperation.NewSuffrageDisjoin(fact)
	if err := op.NodeSign(cmd.Privatekey, cmd.NetworkID.NetworkID(), cmd.node); err != nil {
		return isaacoperation.SuffrageDisjoin{}, errors.Wrap(err, "failed to create suffrage-disjoin operation")
	}

	return op, nil
}
