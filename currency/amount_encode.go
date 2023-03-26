package currency

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (am *Amount) unpack(enc encoder.Encoder, cid string, big string) error {
	e := util.StringErrorFunc("failed to unmarshal Account")

	am.cid = CurrencyID(cid)

	if b, err := NewBigFromString(big); err != nil {
		return e(err, "")
	} else {
		am.big = b
	}

	return nil
}
