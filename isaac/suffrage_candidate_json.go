package isaacoperation

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type suffrageCandidateFactJSONMarshaler struct {
	Address   base.Address   `json:"address"`
	Publickey base.Publickey `json:"publickey"`
	base.BaseFactJSONMarshaler
}

func (fact SuffrageCandidateFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(suffrageCandidateFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Address:               fact.address,
		Publickey:             fact.publickey,
	})
}

type suffrageCandidateFactJSONUnmarshaler struct {
	Address   string `json:"address"`
	Publickey string `json:"publickey"`
	base.BaseFactJSONUnmarshaler
}

func (fact *SuffrageCandidateFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode SuffrageCandidateFact")

	var u suffrageCandidateFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, u.Address, u.Publickey)
}
