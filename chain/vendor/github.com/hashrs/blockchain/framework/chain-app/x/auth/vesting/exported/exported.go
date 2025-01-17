package exported

import (
	"time"

	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	authexported "github.com/hashrs/blockchain/framework/chain-app/x/auth/exported"
)

// VestingAccount defines an account type that vests coins via a vesting schedule.
type VestingAccount interface {
	authexported.Account

	// Delegation and undelegation accounting that returns the resulting base
	// coins amount.
	TrackDelegation(blockTime time.Time, amount sdk.Coins)
	TrackUndelegation(amount sdk.Coins)

	GetVestedCoins(blockTime time.Time) sdk.Coins
	GetVestingCoins(blockTime time.Time) sdk.Coins

	GetStartTime() int64
	GetEndTime() int64

	GetOriginalVesting() sdk.Coins
	GetDelegatedFree() sdk.Coins
	GetDelegatedVesting() sdk.Coins
}
