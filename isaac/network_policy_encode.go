package isaacoperation

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (p *NetworkPolicy) unpack(
	enc encoder.Encoder,
	suffrageCandidateLimiterRule []byte,
	maxOperationsInProposal uint64,
	suffrageCandidateLifespan base.Height,
	maxSuffrageSize uint64,
	suffrageWithdrawLifespan base.Height,
) error {
	e := util.StringErrorFunc("failed to unmarshal NetworkPolicy")

	if err := encoder.Decode(enc, suffrageCandidateLimiterRule, &p.suffrageCandidateLimiterRule); err != nil {
		return e(err, "")
	}

	p.maxOperationsInProposal = maxOperationsInProposal
	p.suffrageCandidateLifespan = suffrageCandidateLifespan
	p.maxSuffrageSize = maxSuffrageSize
	p.suffrageWithdrawLifespan = suffrageWithdrawLifespan

	return nil
}
