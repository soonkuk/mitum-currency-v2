package bsonenc

import (
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

type BSONDecodable interface {
	DecodeBSON([]byte, *Encoder) error
}

func Marshal(v interface{}) ([]byte, error) {
	return bson.Marshal(v)
}

func Unmarshal(b []byte, v interface{}) error {
	return bson.Unmarshal(b, v)
}

type HintedHead struct {
	H string `bson:"_hint"`
}

func NewHintedDoc(h hint.Hint) bson.M {
	return bson.M{"_hint": h.String()}
}

func MergeBSONM(a bson.M, b ...bson.M) bson.M {
	for _, c := range b {
		for k, v := range c {
			a[k] = v
		}
	}

	return a
}
