package keeper

import (
	"time"

	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	"github.com/hashrs/blockchain/framework/chain-app/x/evidence/internal/types"
)

// MaxEvidenceAge returns the maximum age for submitted evidence.
func (k Keeper) MaxEvidenceAge(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyMaxEvidenceAge, &res)
	return
}

// GetParams returns the total set of evidence parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the evidence parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
