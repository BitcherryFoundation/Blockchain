package simulation

import (
	"math/rand"

	simappparams "github.com/hashrs/blockchain/framework/chain-app/simapp/params"
	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	"github.com/hashrs/blockchain/framework/chain-app/x/distribution/keeper"
	"github.com/hashrs/blockchain/framework/chain-app/x/distribution/types"
	govtypes "github.com/hashrs/blockchain/framework/chain-app/x/gov/types"
	"github.com/hashrs/blockchain/framework/chain-app/x/simulation"
)

// OpWeightSubmitCommunitySpendProposal app params key for community spend proposal
const OpWeightSubmitCommunitySpendProposal = "op_weight_submit_community_spend_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []simulation.WeightedProposalContent {
	return []simulation.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightSubmitCommunitySpendProposal,
			DefaultWeight:      simappparams.DefaultWeightCommunitySpendProposal,
			ContentSimulatorFn: SimulateCommunityPoolSpendProposalContent(k),
		},
	}
}

// SimulateCommunityPoolSpendProposalContent generates random community-pool-spend proposal content
// nolint: funlen
func SimulateCommunityPoolSpendProposalContent(k keeper.Keeper) simulation.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simulation.Account) govtypes.Content {
		simAccount, _ := simulation.RandomAcc(r, accs)

		balance := k.GetFeePool(ctx).CommunityPool
		if balance.Empty() {
			return nil
		}

		denomIndex := r.Intn(len(balance))
		amount, err := simulation.RandPositiveInt(r, balance[denomIndex].Amount.TruncateInt())
		if err != nil {
			return nil
		}

		return types.NewCommunityPoolSpendProposal(
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 100),
			simAccount.Address,
			sdk.NewCoins(sdk.NewCoin(balance[denomIndex].Denom, amount)),
		)
	}
}
