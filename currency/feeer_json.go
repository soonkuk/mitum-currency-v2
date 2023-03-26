package currency

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (fa NilFeeer) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(hint.BaseHinter{
		HT: fa.Hint(),
	})
}

func (fa *NilFeeer) UnmarsahlJSON(b []byte) error {
	e := util.StringErrorFunc("failed to unmarshal json of NilFeeer")

	var ht hint.BaseHinter
	if err := util.UnmarshalJSON(b, &ht); err != nil {
		return e(err, "")
	}

	fa.BaseHinter = ht

	return nil
}

type FixedFeeerJSONMarshaler struct {
	hint.BaseHinter
	Receiver base.Address `json:"receiver"`
	Amount   string       `json:"amount"`
}

func (fa FixedFeeer) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(FixedFeeerJSONMarshaler{
		BaseHinter: fa.BaseHinter,
		Receiver:   fa.receiver,
		Amount:     fa.amount.String(),
	})
}

type FixedFeeerJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Receiver string    `json:"receiver"`
	Amount   string    `json:"amount"`
}

func (fa *FixedFeeer) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of FixedFeeer")

	var ufa FixedFeeerJSONUnmarshaler
	if err := enc.Unmarshal(b, &ufa); err != nil {
		return e(err, "")
	}

	return fa.unpack(enc, ufa.Hint, ufa.Receiver, ufa.Amount)
}

type RatioFeeerJSONMarshaler struct {
	hint.BaseHinter
	Receiver base.Address `json:"receiver"`
	Ratio    float64      `json:"ratio"`
	Min      string       `json:"min"`
	Max      string       `json:"max"`
}

func (fa RatioFeeer) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RatioFeeerJSONMarshaler{
		BaseHinter: fa.BaseHinter,
		Receiver:   fa.receiver,
		Ratio:      fa.ratio,
		Min:        fa.min.String(),
		Max:        fa.max.String(),
	})
}

type RatioFeeerJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Receiver string    `json:"receiver"`
	Ratio    float64   `json:"ratio"`
	Min      string    `json:"min"`
	Max      string    `json:"max"`
}

func (fa *RatioFeeer) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of RatioFeeer")

	var ufa RatioFeeerJSONUnmarshaler
	if err := enc.Unmarshal(b, &ufa); err != nil {
		return e(err, "")
	}

	return fa.unpack(enc, ufa.Hint, ufa.Receiver, ufa.Ratio, ufa.Min, ufa.Max)
}
