package isaacoperation

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	NetworkPolicyHint           = hint.MustNewHint("currency-network-policy-v0.0.1")
	NetworkPolicyStateValueHint = hint.MustNewHint("currency-network-policy-state-value-v0.0.1")

	DefaultMaxOperationsInProposal uint64 = 333
	DefaultMaxSuffrageSize         uint64 = 33

	// NOTE suffrage candidate can be approved within lifespan height; almost 15
	// days(based on 5 second for one block)
	DefaultSuffrageCandidateLifespan base.Height = 1 << 18
	DefaultSuffrageWithdrawLifespan              = base.Height(333) //nolint:gomnd //...
)

type NetworkPolicy struct {
	suffrageCandidateLimiterRule base.SuffrageCandidateLimiterRule
	util.DefaultJSONMarshaled
	hint.BaseHinter
	maxOperationsInProposal   uint64
	suffrageCandidateLifespan base.Height
	maxSuffrageSize           uint64
	suffrageWithdrawLifespan  base.Height
}

func DefaultNetworkPolicy() NetworkPolicy {
	return NetworkPolicy{
		BaseHinter:                   hint.NewBaseHinter(NetworkPolicyHint),
		maxOperationsInProposal:      DefaultMaxOperationsInProposal,
		suffrageCandidateLifespan:    DefaultSuffrageCandidateLifespan,
		suffrageCandidateLimiterRule: isaac.NewFixedSuffrageCandidateLimiterRule(1),
		maxSuffrageSize:              DefaultMaxSuffrageSize,
		suffrageWithdrawLifespan:     DefaultSuffrageWithdrawLifespan,
	}
}

func (p NetworkPolicy) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid NetworkPolicy")

	if err := p.BaseHinter.IsValid(NetworkPolicyHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if p.maxOperationsInProposal < 1 {
		return e.Errorf("under zero maxOperationsInProposal")
	}

	if p.maxSuffrageSize < 1 {
		return e.Errorf("under zero maxSuffrageSize")
	}

	switch err := p.suffrageCandidateLifespan.IsValid(nil); {
	case err != nil:
		return e.Wrapf(err, "invalid SuffrageCandidateLifespan")
	case p.suffrageCandidateLifespan <= base.GenesisHeight:
		return e.Errorf("zero SuffrageCandidateLifespan")
	}

	if p.suffrageCandidateLimiterRule == nil {
		return e.Errorf("empty SuffrageCandidateLimiterRule")
	}

	if err := p.suffrageCandidateLimiterRule.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (p NetworkPolicy) HashBytes() []byte {
	var rule []byte

	if p.suffrageCandidateLimiterRule != nil {
		rule = p.suffrageCandidateLimiterRule.HashBytes()
	}

	return util.ConcatBytesSlice(
		util.Uint64ToBytes(p.maxOperationsInProposal),
		p.suffrageCandidateLifespan.Bytes(),
		util.Uint64ToBytes(p.maxSuffrageSize),
		rule,
		p.suffrageWithdrawLifespan.Bytes(),
	)
}

func (p NetworkPolicy) MaxOperationsInProposal() uint64 {
	return p.maxOperationsInProposal
}

func (p NetworkPolicy) SuffrageCandidateLifespan() base.Height {
	return p.suffrageCandidateLifespan
}

func (p NetworkPolicy) SuffrageCandidateLimiterRule() base.SuffrageCandidateLimiterRule {
	return p.suffrageCandidateLimiterRule
}

func (p NetworkPolicy) MaxSuffrageSize() uint64 {
	return p.maxSuffrageSize
}

func (p NetworkPolicy) SuffrageWithdrawLifespan() base.Height {
	return p.suffrageWithdrawLifespan
}

type NetworkPolicyStateValue struct {
	policy base.NetworkPolicy
	hint.BaseHinter
}

func NewNetworkPolicyStateValue(policy base.NetworkPolicy) NetworkPolicyStateValue {
	return NetworkPolicyStateValue{
		BaseHinter: hint.NewBaseHinter(NetworkPolicyStateValueHint),
		policy:     policy,
	}
}

func (s NetworkPolicyStateValue) HashBytes() []byte {
	return s.policy.HashBytes()
}

func (s NetworkPolicyStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid NetworkPolicyStateValue")

	if err := s.BaseHinter.IsValid(NetworkPolicyStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false, s.policy); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (s NetworkPolicyStateValue) Policy() base.NetworkPolicy {
	return s.policy
}
