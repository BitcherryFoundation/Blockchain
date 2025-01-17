package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/hashrs/blockchain/core/consensus/dpos-pbft/libs/kv"

	"github.com/hashrs/blockchain/framework/chain-app/codec"
	"github.com/hashrs/blockchain/framework/chain-app/x/auth/exported"
	"github.com/hashrs/blockchain/framework/chain-app/x/auth/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding auth type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.AddressStoreKeyPrefix):
		var accA, accB exported.Account
		cdc.MustUnmarshalBinaryBare(kvA.Value, &accA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &accB)
		return fmt.Sprintf("%v\n%v", accA, accB)
	case bytes.Equal(kvA.Key, types.GlobalAccountNumberKey):
		var globalAccNumberA, globalAccNumberB uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &globalAccNumberA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &globalAccNumberB)
		return fmt.Sprintf("GlobalAccNumberA: %d\nGlobalAccNumberB: %d", globalAccNumberA, globalAccNumberB)
	default:
		panic(fmt.Sprintf("invalid account key %X", kvA.Key))
	}
}
