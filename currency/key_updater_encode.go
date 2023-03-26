package currency

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *KeyUpdaterFact) unpack(enc encoder.Encoder, tg string, bks []byte, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal KeyUpdaterFact")

	switch ad, err := base.DecodeAddress(tg, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.target = ad
	}

	if hinter, err := enc.Decode(bks); err != nil {
		return err
	} else if k, ok := hinter.(AccountKeys); !ok {
		return util.ErrWrongType.Errorf("expected AccountKeys, not %T", hinter)
	} else {
		fact.keys = k
	}

	fact.currency = CurrencyID(cid)

	return nil
}
