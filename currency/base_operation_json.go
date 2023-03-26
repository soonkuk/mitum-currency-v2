package currency

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

type BaseOperationJSONMarshaler struct {
	Hash  util.Hash   `json:"hash"`
	Fact  base.Fact   `json:"fact"`
	Signs []base.Sign `json:"signs"`
	hint.BaseHinter
}

func (op BaseOperation) JSONMarshaler() BaseOperationJSONMarshaler {
	return BaseOperationJSONMarshaler{
		BaseHinter: op.BaseHinter,
		Hash:       op.h,
		Fact:       op.fact,
		Signs:      op.signs,
	}
}

func (op BaseOperation) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(op.JSONMarshaler())
}

type BaseOperationJSONUnmarshaler struct {
	Hash  valuehash.HashDecoder `json:"hash"`
	Fact  json.RawMessage       `json:"fact"`
	Signs []json.RawMessage     `json:"signs"`
}

func (op *BaseOperation) decodeJSON(b []byte, enc *jsonenc.Encoder, u *BaseOperationJSONUnmarshaler) error {
	if err := enc.Unmarshal(b, u); err != nil {
		return err
	}

	op.h = u.Hash.Hash()

	if err := encoder.Decode(enc, u.Fact, &op.fact); err != nil {
		return errors.WithMessage(err, "failed to decode fact")
	}

	return nil
}

func (op *BaseOperation) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode BaseOperation")

	var u BaseOperationJSONUnmarshaler

	if err := op.decodeJSON(b, enc, &u); err != nil {
		return e(err, "")
	}

	op.signs = make([]base.Sign, len(u.Signs))

	for i := range u.Signs {
		var ub base.BaseSign
		if err := ub.DecodeJSON(u.Signs[i], enc); err != nil {
			return e(err, "failed to decode sign")
		}

		op.signs[i] = ub
	}

	return nil
}
