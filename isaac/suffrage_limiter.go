package isaacoperation

import (
	"math"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	FixedSuffrageCandidateLimiterRuleHint    = hint.MustNewHint("currency-fixed-suffrage-candidate-limiter-rule-v0.0.1")
	MajoritySuffrageCandidateLimiterRuleHint = hint.MustNewHint("currency-majority-suffrage-candidate-limiter-rule-v0.0.1")
)

type FixedSuffrageCandidateLimiterRule struct {
	hint.BaseHinter
	limit uint64
}

func NewFixedSuffrageCandidateLimiterRule(limit uint64) FixedSuffrageCandidateLimiterRule {
	return FixedSuffrageCandidateLimiterRule{
		BaseHinter: hint.NewBaseHinter(FixedSuffrageCandidateLimiterRuleHint),
		limit:      limit,
	}
}

func NewFixedSuffrageCandidateLimiter(rule FixedSuffrageCandidateLimiterRule) base.SuffrageCandidateLimiter {
	return func() (uint64, error) {
		return rule.limit, nil
	}
}

func (l FixedSuffrageCandidateLimiterRule) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid FixedSuffrageCandidateLimiterRule")

	if err := l.BaseHinter.IsValid(FixedSuffrageCandidateLimiterRuleHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (l FixedSuffrageCandidateLimiterRule) Limit() uint64 {
	return l.limit
}

func (l FixedSuffrageCandidateLimiterRule) HashBytes() []byte {
	return util.ConcatBytesSlice(
		l.Hint().Bytes(),
		util.Uint64ToBytes(l.limit),
	)
}

type MajoritySuffrageCandidateLimiterRule struct {
	hint.BaseHinter
	ratio float64
	min   uint64
	max   uint64 // NOTE max < 1 means nolimit
}

func NewMajoritySuffrageCandidateLimiterRule(ratio float64, min, max uint64) MajoritySuffrageCandidateLimiterRule {
	return MajoritySuffrageCandidateLimiterRule{
		BaseHinter: hint.NewBaseHinter(MajoritySuffrageCandidateLimiterRuleHint),
		ratio:      ratio,
		min:        min,
		max:        max,
	}
}

func NewMajoritySuffrageCandidateLimiter(
	rule MajoritySuffrageCandidateLimiterRule,
	getSuffrage func() (uint64, error),
) base.SuffrageCandidateLimiter {
	return func() (uint64, error) {
		i, err := NewCandidatesOfMajoritySuffrageCandidateLimiterRule(rule.ratio, getSuffrage)
		if err != nil {
			return 0, err
		}

		switch {
		case rule.min > 0 && i < rule.min:
			return rule.min, nil
		case rule.max > 0 && i > rule.max:
			return rule.max, nil
		default:
			return i, nil
		}
	}
}

func (l MajoritySuffrageCandidateLimiterRule) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid MajoritySuffrageCandidateLimiterRule")

	if err := l.BaseHinter.IsValid(MajoritySuffrageCandidateLimiterRuleHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if l.ratio < 0 || l.ratio > 1 {
		return e.Errorf("invalid ratio; should be inside 0 <= %0.2f <= 1", l.ratio)
	}

	return nil
}

func (l MajoritySuffrageCandidateLimiterRule) Ratio() float64 {
	return l.ratio
}

func (l MajoritySuffrageCandidateLimiterRule) Min() uint64 {
	return l.min
}

func (l MajoritySuffrageCandidateLimiterRule) Max() uint64 {
	return l.max
}

func (l MajoritySuffrageCandidateLimiterRule) HashBytes() []byte {
	return l.Hint().Bytes()
}

// NewCandidatesOfMajoritySuffrageCandidateLimiterRule find the number of new
// candidates to prevent the current suffrage majority.
func NewCandidatesOfMajoritySuffrageCandidateLimiterRule(
	ratio float64,
	getSuffrage func() (uint64, error),
) (uint64, error) {
	if ratio < 0 || ratio > 1 {
		return 0, errors.Errorf("invalid ratio; should be inside 0 <= %0.2f <= 1", ratio)
	}

	s, err := getSuffrage()
	if err != nil {
		return 0, errors.WithMessage(err, "failed to get the number of candiates for majority suffrage limiter")
	}

	if s < 4 { //nolint:gomnd //...
		return 1, nil
	}

	var nm uint64

	switch f := uint64(math.Floor(float64(s-1) / 3)); {
	case f < 2: //nolint:gomnd //...
		nm = s
	default:
		a := uint64(math.Floor(float64(s-f) * ratio)) // NOTE 20% from fail node

		nm = s - f + a
	}

	ns := nm + uint64(math.Floor(float64(nm-1)/2))

	return ns - s, nil
}
