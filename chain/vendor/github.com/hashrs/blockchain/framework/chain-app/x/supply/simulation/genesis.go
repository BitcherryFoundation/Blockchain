package simulation

// DONTCOVER

import (
	"fmt"

	"github.com/hashrs/blockchain/framework/chain-app/codec"

	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	"github.com/hashrs/blockchain/framework/chain-app/types/module"
	"github.com/hashrs/blockchain/framework/chain-app/x/supply/internal/types"
)

// RandomizedGenState generates a random GenesisState for supply
func RandomizedGenState(simState *module.SimulationState) {
	numAccs := int64(len(simState.Accounts))
	totalSupply := sdk.NewInt(simState.InitialStake * (numAccs + simState.NumBonded))
	supplyGenesis := types.NewGenesisState(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, totalSupply)))

	fmt.Printf("Generated supply parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, supplyGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(supplyGenesis)
}
