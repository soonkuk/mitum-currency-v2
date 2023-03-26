package isaac

import (
	"time"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/localtime"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var ManifestHint = hint.MustNewHint("digest-manifest-v0.0.1")

type Manifest struct {
	proposedAt     time.Time
	statesTree     util.Hash
	h              util.Hash
	previous       util.Hash
	proposal       util.Hash
	operationsTree util.Hash
	suffrage       util.Hash
	height         base.Height
	hint.BaseHinter
}

func NewManifest(
	height base.Height,
	previous,
	proposal,
	operationsTree,
	statesTree,
	suffrage util.Hash,
	proposedAt time.Time,
) Manifest {
	m := Manifest{
		BaseHinter:     hint.NewBaseHinter(ManifestHint),
		height:         height,
		previous:       previous,
		proposal:       proposal,
		operationsTree: operationsTree,
		statesTree:     statesTree,
		suffrage:       suffrage,
		proposedAt:     proposedAt,
	}

	m.h = m.generateHash()

	return m
}

func (m Manifest) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid manifest")

	if err := m.BaseHinter.IsValid(ManifestHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false,
		m.height,
		m.proposal,
		util.DummyIsValider(func([]byte) error {
			if m.proposedAt.IsZero() {
				return util.ErrInvalid.Errorf("empty proposedAt")
			}

			return nil
		}),
	); err != nil {
		return e.Wrap(err)
	}

	if m.height != base.GenesisHeight {
		if err := util.CheckIsValiders(nil, false, m.previous); err != nil {
			return e.Wrap(err)
		}
	}

	if err := util.CheckIsValiders(nil, true,
		m.operationsTree,
		m.statesTree,
		m.suffrage,
	); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (m Manifest) Hash() util.Hash {
	return m.h
}

func (m Manifest) Height() base.Height {
	return m.height
}

func (m Manifest) Previous() util.Hash {
	return m.previous
}

func (m Manifest) Proposal() util.Hash {
	return m.proposal
}

func (m Manifest) OperationsTree() util.Hash {
	return m.operationsTree
}

func (m Manifest) StatesTree() util.Hash {
	return m.statesTree
}

func (m Manifest) Suffrage() util.Hash {
	return m.suffrage
}

func (m Manifest) ProposedAt() time.Time {
	return m.proposedAt
}

func (m Manifest) generateHash() util.Hash {
	return valuehash.NewSHA256(util.ConcatByters(
		m.height,
		m.previous,
		m.proposal,
		m.operationsTree,
		m.statesTree,
		m.suffrage,
		localtime.New(m.proposedAt),
	))
}
