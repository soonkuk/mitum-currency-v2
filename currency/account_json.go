package currency

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

type AccountJSONMarshaler struct {
	hint.BaseHinter
	Hash    util.Hash    `json:"hash"`
	Address base.Address `json:"address"`
	Keys    AccountKeys  `json:"keys"`
}

func (ac Account) EncodeJSON() AccountJSONMarshaler {
	return AccountJSONMarshaler{
		BaseHinter: ac.BaseHinter,
		Hash:       ac.h,
		Address:    ac.address,
		Keys:       ac.keys,
	}
}

func (ac Account) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AccountJSONMarshaler{
		BaseHinter: ac.BaseHinter,
		Hash:       ac.h,
		Address:    ac.address,
		Keys:       ac.keys,
	})
}

type AccountJSONUnmarshaler struct {
	Hint    hint.Hint             `json:"_hint"`
	Hash    valuehash.HashDecoder `json:"hash"`
	Address string                `json:"address"`
	Keys    json.RawMessage       `json:"keys"`
}

func (ac *Account) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Account")

	var uac AccountJSONUnmarshaler
	if err := enc.Unmarshal(b, &uac); err != nil {
		return e(err, "")
	}

	return ac.unpack(enc, uac.Hint, uac.Hash, uac.Address, uac.Keys)
}
