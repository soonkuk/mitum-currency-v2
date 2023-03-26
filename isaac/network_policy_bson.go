package isaacoperation

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/digest/util/bson"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (p NetworkPolicy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":                       p.Hint().String(),
			"suffrage_candidate_limiter":  p.suffrageCandidateLimiterRule,
			"max_operations_in_proposal":  p.maxOperationsInProposal,
			"suffrage_candidate_lifespan": p.suffrageCandidateLifespan,
			"max_suffrage_size":           p.maxSuffrageSize,
			"suffrage_withdraw_lifespan":  p.suffrageWithdrawLifespan,
		},
	)
}

type NetworkPolicyBSONUnMarshaler struct {
	Hint                         string      `bson:"_hint"`
	SuffrageCandidateLimiterRule bson.Raw    `bson:"suffrage_candidate_limiter"`
	MaxOperationsInProposal      uint64      `bson:"max_operations_in_proposal"`
	SuffrageCandidateLifespan    base.Height `bson:"suffrage_candidate_lifespan"`
	MaxSuffrageSize              uint64      `bson:"max_suffrage_size"`
	SuffrageWithdrawLifespan     base.Height `bson:"suffrage_withdraw_lifespan"`
}

func (p *NetworkPolicy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of NetworkPolicy")

	var u NetworkPolicyBSONUnMarshaler
	if err := bson.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	p.BaseHinter = hint.NewBaseHinter(ht)

	return p.unpack(enc, u.SuffrageCandidateLimiterRule, u.MaxOperationsInProposal, u.SuffrageCandidateLifespan, u.MaxSuffrageSize, u.SuffrageWithdrawLifespan)
}

func (s NetworkPolicyStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  s.Hint().String(),
			"policy": s.policy,
		},
	)
}

type NetworkPolicyStateValueBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Policy bson.Raw `bson:"policy"`
}

func (s *NetworkPolicyStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode NetworkPolicyStateValue")

	var u NetworkPolicyStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(ht)

	if err := encoder.Decode(enc, u.Policy, &s.policy); err != nil {
		return e(err, "")
	}

	return nil
}
