package isaacoperation

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type suffrageJoinFactJSONMarshaler struct {
	Candidate base.Address `json:"candidate"`
	base.BaseFactJSONMarshaler
	Start base.Height `json:"start_height"`
}

func (fact SuffrageJoinFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(suffrageJoinFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Candidate:             fact.candidate,
		Start:                 fact.start,
	})
}

type suffrageJoinFactJSONUnmarshaler struct {
	Candidate string `json:"candidate"`
	base.BaseFactJSONUnmarshaler
	Start base.Height `json:"start_height"`
}

func (fact *SuffrageJoinFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode SuffrageJoinFact")

	var u suffrageJoinFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, u.Candidate, u.Start)
}

type suffrageGenesisJoinFactJSONMarshaler struct {
	Nodes []base.Node `json:"nodes"`
	base.BaseFactJSONMarshaler
}

func (fact SuffrageGenesisJoinFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(suffrageGenesisJoinFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Nodes:                 fact.nodes,
	})
}

type suffrageGenesisJoinFactJSONUnmarshaler struct {
	Nodes []json.RawMessage `json:"nodes"`
	base.BaseFactJSONUnmarshaler
}

func (fact *SuffrageGenesisJoinFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode SuffrageGenesisJoinFact")

	var u suffrageGenesisJoinFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	fact.nodes = make([]base.Node, len(u.Nodes))

	for i := range u.Nodes {
		if err := encoder.Decode(enc, u.Nodes[i], &fact.nodes[i]); err != nil {
			return e(err, "")
		}
	}

	return nil
}
