package keeper

import (
	"fmt"

	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	sdkerrors "github.com/hashrs/blockchain/framework/chain-app/types/errors"
	"github.com/hashrs/blockchain/framework/chain-app/x/distribution/types"
)

// HandleCommunityPoolSpendProposal is a handler for executing a passed community spend proposal
func HandleCommunityPoolSpendProposal(ctx sdk.Context, k Keeper, p types.CommunityPoolSpendProposal) error {
	if k.blacklistedAddrs[p.Recipient.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is blacklisted from receiving external funds", p.Recipient)
	}

	err := k.DistributeFromFeePool(ctx, p.Amount, p.Recipient)
	if err != nil {
		return err
	}

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("transferred %s from the community pool to recipient %s", p.Amount, p.Recipient))
	return nil
}
