package keeper

import (
	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	"github.com/hashrs/blockchain/framework/chain-app/x/slashing/internal/types"
)

// Unjail calls the staking Unjail function to unjail a validator if the
// jailed period has concluded
func (k Keeper) Unjail(ctx sdk.Context, validatorAddr sdk.ValAddress) error {
	validator := k.sk.Validator(ctx, validatorAddr)
	if validator == nil {
		return types.ErrNoValidatorForAddress
	}

	// cannot be unjailed if no self-delegation exists
	selfDel := k.sk.Delegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr)
	if selfDel == nil {
		return types.ErrMissingSelfDelegation
	}

	if validator.TokensFromShares(selfDel.GetShares()).TruncateInt().LT(validator.GetMinSelfDelegation()) {
		return types.ErrSelfDelegationTooLowToUnjail
	}

	// cannot be unjailed if not jailed
	if !validator.IsJailed() {
		return types.ErrValidatorNotJailed
	}

	consAddr := sdk.ConsAddress(validator.GetConsPubKey().Address())

	info, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		return types.ErrNoValidatorForAddress
	}

	// cannot be unjailed if tombstoned
	if info.Tombstoned {
		return types.ErrValidatorJailed
	}

	// cannot be unjailed until out of jail
	if ctx.BlockHeader().Time.Before(info.JailedUntil) {
		return types.ErrValidatorJailed
	}

	k.sk.Unjail(ctx, consAddr)
	return nil
}
