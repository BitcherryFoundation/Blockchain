package types

import (
	"time"

	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	stakingexported "github.com/hashrs/blockchain/framework/chain-app/x/staking/exported"

	"github.com/hashrs/blockchain/core/consensus/dpos-pbft/crypto"
)

type (
	// StakingKeeper defines the staking module interface contract needed by the
	// evidence module.
	StakingKeeper interface {
		ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingexported.ValidatorI
	}

	// SlashingKeeper defines the slashing module interface contract needed by the
	// evidence module.
	SlashingKeeper interface {
		GetPubkey(sdk.Context, crypto.Address) (crypto.PubKey, error)
		IsTombstoned(sdk.Context, sdk.ConsAddress) bool
		HasValidatorSigningInfo(sdk.Context, sdk.ConsAddress) bool
		Tombstone(sdk.Context, sdk.ConsAddress)
		Slash(sdk.Context, sdk.ConsAddress, sdk.Dec, int64, int64)
		SlashFractionDoubleSign(sdk.Context) sdk.Dec
		Jail(sdk.Context, sdk.ConsAddress)
		JailUntil(sdk.Context, sdk.ConsAddress, time.Time)
	}
)
