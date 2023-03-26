package currency

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/digest/util/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

func (ca Address) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, bsoncore.AppendString(nil, ca.String()), nil
}

func (ca *Address) DecodeBSON(b []byte, _ *bsonenc.Encoder) error {
	*ca = NewAddress(string(b))

	return nil
}

func (ca EthAddress) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, bsoncore.AppendString(nil, ca.String()), nil
}

func (ca *EthAddress) DecodeBSON(b []byte, _ *bsonenc.Encoder) error {
	*ca = NewEthAddress(string(b))

	return nil
}
