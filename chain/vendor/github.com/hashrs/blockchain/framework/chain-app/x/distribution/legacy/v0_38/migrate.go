package v0_38

// DONTCOVER
// nolint

import (
	v036distr "github.com/hashrs/blockchain/framework/chain-app/x/distribution/legacy/v0_36"
)

// Migrate accepts exported genesis state from v0.36 or v0.37 and migrates it to
// v0.38 genesis state. All entries are identical except for parameters.
func Migrate(oldGenState v036distr.GenesisState) GenesisState {
	params := Params{
		CommunityTax:        oldGenState.CommunityTax,
		BaseProposerReward:  oldGenState.BaseProposerReward,
		BonusProposerReward: oldGenState.BonusProposerReward,
		WithdrawAddrEnabled: oldGenState.WithdrawAddrEnabled,
	}

	return NewGenesisState(
		params, oldGenState.FeePool,
		oldGenState.DelegatorWithdrawInfos, oldGenState.PreviousProposer,
		oldGenState.OutstandingRewards, oldGenState.ValidatorAccumulatedCommissions,
		oldGenState.ValidatorHistoricalRewards, oldGenState.ValidatorCurrentRewards,
		oldGenState.DelegatorStartingInfos, oldGenState.ValidatorSlashEvents,
	)
}
