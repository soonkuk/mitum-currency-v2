package currency

import (
	"bytes"
	"sort"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/util"
)

var (
	AccountKeyHint  = hint.MustNewHint("mitum-currency-key-v0.0.1")
	AccountKeysHint = hint.MustNewHint("mitum-currency-keys-v0.0.1")
)

var MaxAccountKeyInKeys = 10

type AccountKey interface {
	hint.Hinter
	util.IsValider
	util.Byter
	Key() base.Publickey
	Weight() uint
	Equal(AccountKey) bool
}

type AccountKeys interface {
	hint.Hinter
	util.IsValider
	util.Byter
	util.Hasher
	Threshold() uint
	Keys() []AccountKey
	Key(base.Publickey) (AccountKey, bool)
	Equal(AccountKeys) bool
}

type BaseAccountKey struct {
	hint.BaseHinter
	k base.Publickey
	w uint
}

func NewBaseAccountKey(k base.Publickey, w uint) (BaseAccountKey, error) {
	ky := BaseAccountKey{BaseHinter: hint.NewBaseHinter(AccountKeyHint), k: k, w: w}

	return ky, ky.IsValid(nil)
}

func (ky BaseAccountKey) IsValid([]byte) error {
	if ky.w < 1 || ky.w > 100 {
		return util.ErrInvalid.Errorf("invalid key weight, 1 <= weight <= 100")
	}

	return util.CheckIsValiders(nil, false, ky.k)
}

func (ky BaseAccountKey) Weight() uint {
	return ky.w
}

func (ky BaseAccountKey) Key() base.Publickey {
	return ky.k
}

func (ky BaseAccountKey) Bytes() []byte {
	return util.ConcatBytesSlice(ky.k.Bytes(), util.UintToBytes(ky.w))
}

func (ky BaseAccountKey) Equal(b AccountKey) bool {
	if ky.w != b.Weight() {
		return false
	}

	if !ky.k.Equal(b.Key()) {
		return false
	}

	return true
}

type BaseAccountKeys struct {
	hint.BaseHinter
	h         util.Hash
	keys      []AccountKey
	threshold uint
}

func EmptyBaseAccountKeys() BaseAccountKeys {
	return BaseAccountKeys{BaseHinter: hint.NewBaseHinter(AccountKeysHint)}
}

func NewBaseAccountKeys(keys []AccountKey, threshold uint) (BaseAccountKeys, error) {
	ks := BaseAccountKeys{BaseHinter: hint.NewBaseHinter(AccountKeysHint), keys: keys, threshold: threshold}
	h, err := ks.GenerateHash()
	if err != nil {
		return BaseAccountKeys{}, err
	}
	ks.h = h

	return ks, ks.IsValid(nil)
}

func (ks BaseAccountKeys) Hash() util.Hash {
	return ks.h
}

func (ks BaseAccountKeys) GenerateHash() (util.Hash, error) {
	return valuehash.NewSHA256(ks.Bytes()), nil
}

func (ks BaseAccountKeys) GenerateKeccakHash() (util.Hash, error) {
	var b valuehash.L32
	h := crypto.Keccak256(ks.Bytes()[:])
	copy(b[:], h[:])
	return b, nil
}

func (ks BaseAccountKeys) Bytes() []byte {
	bs := make([][]byte, len(ks.keys)+1)

	// NOTE sorted by Key.Key()
	sort.Slice(ks.keys, func(i, j int) bool {
		return bytes.Compare(ks.keys[i].Key().Bytes(), ks.keys[j].Key().Bytes()) < 0
	})
	for i := range ks.keys {
		bs[i] = ks.keys[i].Bytes()
	}

	bs[len(ks.keys)] = util.UintToBytes(ks.threshold)

	return util.ConcatBytesSlice(bs...)
}

func (ks BaseAccountKeys) IsValid([]byte) error {
	if ks.threshold < 1 || ks.threshold > 100 {
		return util.ErrInvalid.Errorf("invalid threshold, %d, should be 1 <= threshold <= 100", ks.threshold)
	}

	if err := util.CheckIsValiders(nil, false, ks.h); err != nil {
		return err
	}

	if n := len(ks.keys); n < 1 {
		return util.ErrInvalid.Errorf("empty keys")
	} else if n > MaxAccountKeyInKeys {
		return util.ErrInvalid.Errorf("keys over %d, %d", MaxAccountKeyInKeys, n)
	}

	m := map[string]struct{}{}
	for i := range ks.keys {
		k := ks.keys[i]
		if err := util.CheckIsValiders(nil, false, k); err != nil {
			return err
		}

		if _, found := m[k.Key().String()]; found {
			return util.ErrInvalid.Errorf("duplicated keys found")
		}

		m[k.Key().String()] = struct{}{}
	}

	var totalWeight uint
	for i := range ks.keys {
		totalWeight += ks.keys[i].Weight()
	}

	if totalWeight < ks.threshold {
		return util.ErrInvalid.Errorf("sum of weight under threshold, %d < %d", totalWeight, ks.threshold)
	}

	if h, err := ks.GenerateHash(); err != nil {
		return err
	} else if !ks.h.Equal(h) {
		return util.ErrInvalid.Errorf("hash not matched")
	}

	return nil
}

func (ks BaseAccountKeys) Threshold() uint {
	return ks.threshold
}

func (ks BaseAccountKeys) Keys() []AccountKey {
	return ks.keys
}

func (ks BaseAccountKeys) Key(k base.Publickey) (AccountKey, bool) {
	for i := range ks.keys {
		ky := ks.keys[i]
		if ky.Key().Equal(k) {
			return ky, true
		}
	}

	return BaseAccountKey{}, false
}

func (ks BaseAccountKeys) Equal(b AccountKeys) bool {
	if ks.threshold != b.Threshold() {
		return false
	}

	if len(ks.keys) != len(b.Keys()) {
		return false
	}

	sort.Slice(ks.keys, func(i, j int) bool {
		return bytes.Compare(ks.keys[i].Key().Bytes(), ks.keys[j].Key().Bytes()) < 0
	})

	bkeys := b.Keys()
	sort.Slice(bkeys, func(i, j int) bool {
		return bytes.Compare(bkeys[i].Key().Bytes(), bkeys[j].Key().Bytes()) < 0
	})

	for i := range ks.keys {
		if !ks.keys[i].Equal(bkeys[i]) {
			return false
		}
	}

	return true
}

func checkThreshold(fs []base.Sign, keys AccountKeys) error {
	var sum uint
	for i := range fs {
		ky, found := keys.Key(fs[i].Signer())
		if !found {
			return errors.Errorf("unknown key found, %s", fs[i].Signer())
		}
		sum += ky.Weight()
	}

	if sum < keys.Threshold() {
		return errors.Errorf("not passed threshold, sum=%d < threshold=%d", sum, keys.Threshold())
	}

	return nil
}
