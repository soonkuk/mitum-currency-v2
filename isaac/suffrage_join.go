package isaacoperation

import (
	"bytes"
	"context"

	"github.com/ProtoconNet/mitum-currency/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	SuffrageJoinFactHint        = hint.MustNewHint("currency-suffrage-join-fact-v0.0.1")
	SuffrageJoinHint            = hint.MustNewHint("currency-suffrage-join-operation-v0.0.1")
	SuffrageGenesisJoinFactHint = hint.MustNewHint("currency-suffrage-genesis-join-fact-v0.0.1")
	SuffrageGenesisJoinHint     = hint.MustNewHint("currency-suffrage-genesis-join-operation-v0.0.1")
)

type SuffrageJoinFact struct {
	candidate base.Address
	base.BaseFact
	start base.Height
}

func NewSuffrageJoinFact(
	token base.Token,
	candidate base.Address,
	start base.Height,
) SuffrageJoinFact {
	fact := SuffrageJoinFact{
		BaseFact:  base.NewBaseFact(SuffrageJoinFactHint, token),
		candidate: candidate,
		start:     start,
	}

	fact.SetHash(fact.hash())

	return fact
}

func (fact SuffrageJoinFact) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid SuffrageJoinFact")

	if err := util.CheckIsValiders(nil, false, fact.BaseFact, fact.candidate, fact.start); err != nil {
		return e.Wrap(err)
	}

	if !fact.Hash().Equal(fact.hash()) {
		return e.Errorf("hash does not match")
	}

	return nil
}

func (fact SuffrageJoinFact) Candidate() base.Address {
	return fact.candidate
}

func (fact SuffrageJoinFact) Start() base.Height {
	return fact.start
}

func (fact SuffrageJoinFact) hash() util.Hash {
	return valuehash.NewSHA256(util.ConcatByters(
		util.BytesToByter(fact.Token()),
		fact.candidate,
		fact.start,
	))
}

type SuffrageGenesisJoinFact struct {
	nodes []base.Node
	base.BaseFact
}

func NewSuffrageGenesisJoinFact(
	nodes []base.Node,
	networkID base.NetworkID,
) SuffrageGenesisJoinFact {
	fact := SuffrageGenesisJoinFact{
		BaseFact: base.NewBaseFact(SuffrageGenesisJoinFactHint, base.Token(networkID)),
		nodes:    nodes,
	}

	fact.SetHash(fact.hash())

	return fact
}

func (fact SuffrageGenesisJoinFact) IsValid(networkID []byte) error {
	e := util.ErrInvalid.Errorf("invalid SuffrageGenesisJoinFact")

	if len(fact.nodes) < 1 {
		return e.Errorf("empty nodes")
	}

	vs := make([]util.IsValider, len(fact.nodes)+1)
	vs[0] = fact.BaseFact

	for i := range fact.nodes {
		vs[i+1] = fact.nodes[i]
	}

	if err := util.CheckIsValiders(nil, false, vs...); err != nil {
		return e.Wrap(err)
	}

	if _, found := util.IsDuplicatedSlice(fact.nodes, func(i base.Node) (bool, string) {
		return true, i.Address().String()
	}); found {
		return e.Errorf("duplicated node found")
	}

	if !bytes.Equal(fact.BaseFact.Token(), networkID) {
		return e.Errorf("wrong token")
	}

	if !fact.Hash().Equal(fact.hash()) {
		return e.Errorf("hash does not match")
	}

	return nil
}

func (fact SuffrageGenesisJoinFact) Nodes() []base.Node {
	return fact.nodes
}

func (fact SuffrageGenesisJoinFact) hash() util.Hash {
	return valuehash.NewSHA256(util.ConcatByters(
		util.BytesToByter(fact.Token()),
		util.DummyByter(func() []byte {
			var b bytes.Buffer

			for i := range fact.nodes {
				_, _ = b.Write(fact.nodes[i].HashBytes())
			}

			return b.Bytes()
		}),
	))
}

type SuffrageJoin struct {
	currency.BaseNodeOperation
}

func NewSuffrageJoin(fact SuffrageJoinFact) SuffrageJoin {
	return SuffrageJoin{
		BaseNodeOperation: currency.NewBaseNodeOperation(SuffrageJoinHint, fact),
	}
}

func (op *SuffrageJoin) SetToken(t base.Token) error {
	fact := op.Fact().(SuffrageJoinFact) //nolint:forcetypeassert //...

	if err := fact.SetToken(t); err != nil {
		return err
	}

	fact.SetHash(fact.hash())

	op.BaseNodeOperation.SetFact(fact)

	return nil
}

func (op SuffrageJoin) IsValid(networkID []byte) error {
	e := util.ErrInvalid.Errorf("invalid SuffrageJoin")

	if err := op.BaseNodeOperation.IsValid(networkID); err != nil {
		return e.Wrap(err)
	}

	fact, ok := op.Fact().(SuffrageJoinFact)
	if !ok {
		return e.Errorf("not SuffrageJoinFact, %T", op.Fact())
	}

	var foundsigner bool

	sfs := op.Signs()

	for i := range sfs {
		ns := sfs[i].(base.NodeSign) //nolint:forcetypeassert //...

		if !ns.Node().Equal(fact.Candidate()) {
			continue
		}

		foundsigner = true

		break
	}

	if !foundsigner {
		return e.Errorf("not signed by candidate")
	}

	return nil
}

// SuffrageGenesisJoin is only for used for genesis block
type SuffrageGenesisJoin struct {
	currency.BaseOperation
}

func NewSuffrageGenesisJoin(fact SuffrageGenesisJoinFact) SuffrageGenesisJoin {
	return SuffrageGenesisJoin{
		BaseOperation: currency.NewBaseOperation(SuffrageGenesisJoinHint, fact),
	}
}

func (op SuffrageGenesisJoin) IsValid(networkID []byte) error {
	e := util.ErrInvalid.Errorf("invalid SuffrageGenesisJoin")

	if err := op.BaseOperation.IsValid(networkID); err != nil {
		return e.Wrap(err)
	}

	if len(op.Signs()) > 1 {
		return e.Errorf("multiple signs found")
	}

	if _, ok := op.Fact().(SuffrageGenesisJoinFact); !ok {
		return e.Errorf("not SuffrageGenesisJoinFact, %T", op.Fact())
	}

	return nil
}

func (SuffrageGenesisJoin) PreProcess(
	ctx context.Context, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	switch _, found, err := getStateFunc(isaac.SuffrageStateKey); {
	case err != nil:
		return ctx, base.NewBaseOperationProcessReasonError("failed to check suffrage state: %w", err), nil
	case found:
		return ctx, base.NewBaseOperationProcessReasonError("suffrage state already exists"), nil
	default:
		return ctx, nil, nil
	}
}

func (op SuffrageGenesisJoin) Process(context.Context, base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact := op.Fact().(SuffrageGenesisJoinFact) //nolint:forcetypeassert //...

	fnodes := fact.Nodes()
	nodes := make([]base.SuffrageNodeStateValue, len(fnodes))

	for i := range fnodes {
		nodes[i] = isaac.NewSuffrageNodeStateValue(fnodes[i], base.GenesisHeight+1)
	}

	return []base.StateMergeValue{
		currency.NewBaseStateMergeValue(
			isaac.SuffrageStateKey,
			isaac.NewSuffrageNodesStateValue(base.GenesisHeight, nodes),
			nil,
		),
	}, nil, nil
}
