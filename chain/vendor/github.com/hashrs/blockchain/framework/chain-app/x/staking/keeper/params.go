package keeper

import (
	"time"

	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	"github.com/hashrs/blockchain/framework/chain-app/x/params"
	"github.com/hashrs/blockchain/framework/chain-app/x/staking/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// UnbondingTime
func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
	return
}

// MaxValidators - Maximum number of validators
func (k Keeper) MaxValidators(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

// MaxEntries - Maximum number of simultaneous unbonding
// delegations or redelegations (per pair/trio)
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyMaxEntries, &res)
	return
}

// HistoricalEntries = number of historical info entries
// to persist in store
func (k Keeper) HistoricalEntries(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyHistoricalEntries, &res)
	return
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.UnbondingTime(ctx),
		k.MaxValidators(ctx),
		k.MaxEntries(ctx),
		k.HistoricalEntries(ctx),
		k.BondDenom(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
