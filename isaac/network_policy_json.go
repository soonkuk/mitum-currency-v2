package isaacoperation

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type networkPolicyJSONMarshaler struct {
	// revive:disable-next-line:line-length-limit
	SuffrageCandidateLimiterRule base.SuffrageCandidateLimiterRule `json:"suffrage_candidate_limiter"` //nolint:tagliatelle //...
	hint.BaseHinter
	MaxOperationsInProposal   uint64      `json:"max_operations_in_proposal"`
	SuffrageCandidateLifespan base.Height `json:"suffrage_candidate_lifespan"`
	MaxSuffrageSize           uint64      `json:"max_suffrage_size"`
	SuffrageWithdrawLifespan  base.Height `json:"suffrage_withdraw_lifespan"`
}

func (p NetworkPolicy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(networkPolicyJSONMarshaler{
		BaseHinter:                   p.BaseHinter,
		MaxOperationsInProposal:      p.maxOperationsInProposal,
		SuffrageCandidateLifespan:    p.suffrageCandidateLifespan,
		SuffrageCandidateLimiterRule: p.suffrageCandidateLimiterRule,
		MaxSuffrageSize:              p.maxSuffrageSize,
		SuffrageWithdrawLifespan:     p.suffrageWithdrawLifespan,
	})
}

type networkPolicyJSONUnmarshaler struct {
	SuffrageCandidateLimiterRule json.RawMessage `json:"suffrage_candidate_limiter"` //nolint:tagliatelle //...
	MaxOperationsInProposal      uint64          `json:"max_operations_in_proposal"`
	SuffrageCandidateLifespan    base.Height     `json:"suffrage_candidate_lifespan"`
	MaxSuffrageSize              uint64          `json:"max_suffrage_size"`
	SuffrageWithdrawLifespan     base.Height     `json:"suffrage_withdraw_lifespan"`
}

func (p *NetworkPolicy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to unmarshal NetworkPolicy")

	var u networkPolicyJSONUnmarshaler
	if err := util.UnmarshalJSON(b, &u); err != nil {
		return e(err, "")
	}

	return p.unpack(enc, u.SuffrageCandidateLimiterRule, u.MaxOperationsInProposal, u.SuffrageCandidateLifespan, u.MaxSuffrageSize, u.SuffrageWithdrawLifespan)
}

type NetworkPolicyStateValueJSONMarshaler struct {
	Policy base.NetworkPolicy `json:"policy"`
	hint.BaseHinter
}

func (s NetworkPolicyStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(NetworkPolicyStateValueJSONMarshaler{
		BaseHinter: s.BaseHinter,
		Policy:     s.policy,
	})
}

type NetworkPolicyStateValueJSONUnmarshaler struct {
	Policy json.RawMessage `json:"policy"`
}

func (s *NetworkPolicyStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode NetworkPolicyStateValue")

	var u NetworkPolicyStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	if err := encoder.Decode(enc, u.Policy, &s.policy); err != nil {
		return e(err, "")
	}

	return nil
}
