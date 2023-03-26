package isaacoperation

import (
	"github.com/ProtoconNet/mitum-currency/currency"
	bsonenc "github.com/ProtoconNet/mitum-currency/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

func (fact GenesisNetworkPolicyFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  fact.Hint().String(),
			"policy": fact.policy,
			"hash":   fact.BaseFact.Hash().String(),
			"token":  fact.BaseFact.Token(),
		},
	)
}

type GenesisNetworkPolicyFactBSONUnMarshaler struct {
	Hint   string   `bson:"_hint"`
	Policy bson.Raw `bson:"policy"`
}

func (fact *GenesisNetworkPolicyFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of GenesisNetworkPolicyFact")

	var u currency.BaseFactBSONUnmarshaler

	err := enc.Unmarshal(b, &u)
	if err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(u.Hash))
	fact.BaseFact.SetToken(u.Token)

	var uf GenesisNetworkPolicyFactBSONUnMarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	if err := encoder.Decode(enc, uf.Policy, &fact.policy); err != nil {
		return e(err, "")
	}

	return nil
}
