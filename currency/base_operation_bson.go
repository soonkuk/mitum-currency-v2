package currency

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/digest/util/bson"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

type BaseOperationBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	Hash string   `bson:"hash"`
	Fact bson.Raw `bson:"fact"`
	// Signs []bson.Raw      `bson:"signs"`
}

func (op BaseOperation) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		},
	)
}

func (op *BaseOperation) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of BaseOperation")

	var u BaseOperationBSONUnmarshaler

	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	op.BaseHinter = hint.NewBaseHinter(ht)
	op.h = valuehash.NewBytesFromString(u.Hash)

	var fact base.Fact
	if err := encoder.Decode(enc, u.Fact, &fact); err != nil {
		return e(err, "failed to decode fact")
	}

	op.SetFact(fact)

	return nil
}

func (op *BaseNodeOperation) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of BaseNodeOperation")

	var u BaseOperationBSONUnmarshaler

	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	op.BaseOperation.BaseHinter = hint.NewBaseHinter(ht)
	op.BaseOperation.h = valuehash.NewBytesFromString(u.Hash)

	var fact base.Fact
	if err := encoder.Decode(enc, u.Fact, &fact); err != nil {
		return e(err, "failed to decode fact")
	}

	op.BaseOperation.SetFact(fact)

	return nil
}
