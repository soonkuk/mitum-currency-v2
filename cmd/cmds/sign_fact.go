package cmds

import (
	"fmt"
	"os"

	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base/operation"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/logging"
)

type SignFactCommand struct {
	Privatekey PrivatekeyFlag `arg:"" name:"privatekey" help:"sender's privatekey" required:""`
	NetworkID  string         `name:"network-id" help:"network-id" required:""`
	Pretty     bool           `name:"pretty" help:"pretty format"`
	Seal       string         `help:"seal" optional:"" type:"existingfile"`
}

func (cmd *SignFactCommand) Run(log logging.Logger) error {
	var sl operation.Seal
	if s, err := loadSealFromInput(cmd.Seal); err != nil {
		return err
	} else if so, ok := s.(operation.Seal); !ok {
		return xerrors.Errorf("seal is not operation.Seal, %T", s)
	} else if _, ok := so.(operation.SealUpdater); !ok {
		return xerrors.Errorf("seal is not operation.SealUpdater, %T", so)
	} else if err := so.IsValid([]byte(cmd.NetworkID)); err != nil {
		return xerrors.Errorf("invalid seal: %w", err)
	} else {
		sl = so
	}
	log.Debug().Hinted("seal", sl.Hash()).Msg("seal loaded")

	nops := make([]operation.Operation, len(sl.Operations()))
	for i := range sl.Operations() {
		op := sl.Operations()[i]

		var fsu operation.FactSignUpdater
		if u, ok := op.(operation.FactSignUpdater); !ok {
			log.Debug().
				Interface("operation", op).
				Hinted("operation_type", op.Hint()).
				Msg("not operation.FactSignUpdater")

			nops[i] = op
		} else {
			fsu = u
		}

		if sig, err := operation.NewFactSignature(cmd.Privatekey, op.Fact(), []byte(cmd.NetworkID)); err != nil {
			return err
		} else {
			f := operation.NewBaseFactSign(cmd.Privatekey.Publickey(), sig)

			if nop, err := fsu.AddFactSigns(f); err != nil {
				return err
			} else {
				nops[i] = nop.(operation.Operation)
			}
		}
	}

	sl = sl.(operation.SealUpdater).SetOperations(nops).(operation.Seal)

	if s, err := signSeal(sl, cmd.Privatekey, []byte(cmd.NetworkID)); err != nil {
		return err
	} else {
		sl = s.(operation.Seal)

		log.Debug().Msg("seal signed")
	}

	_, _ = fmt.Fprintln(os.Stdout, string(jsonenc.MustMarshalIndent(sl)))

	return nil
}
