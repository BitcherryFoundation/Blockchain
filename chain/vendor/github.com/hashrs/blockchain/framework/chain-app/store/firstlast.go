package store

import (
	"bytes"

	tmkv "github.com/hashrs/blockchain/core/consensus/dpos-pbft/libs/kv"

	"github.com/hashrs/blockchain/framework/chain-app/store/types"
)

// Gets the first item.
func First(st KVStore, start, end []byte) (kv tmkv.Pair, ok bool) {
	iter := st.Iterator(start, end)
	if !iter.Valid() {
		return kv, false
	}
	defer iter.Close()

	return tmkv.Pair{Key: iter.Key(), Value: iter.Value()}, true
}

// Gets the last item.  `end` is exclusive.
func Last(st KVStore, start, end []byte) (kv tmkv.Pair, ok bool) {
	iter := st.ReverseIterator(end, start)
	if !iter.Valid() {
		if v := st.Get(start); v != nil {
			return tmkv.Pair{Key: types.Cp(start), Value: types.Cp(v)}, true
		}
		return kv, false
	}
	defer iter.Close()

	if bytes.Equal(iter.Key(), end) {
		// Skip this one, end is exclusive.
		iter.Next()
		if !iter.Valid() {
			return kv, false
		}
	}

	return tmkv.Pair{Key: iter.Key(), Value: iter.Value()}, true
}
