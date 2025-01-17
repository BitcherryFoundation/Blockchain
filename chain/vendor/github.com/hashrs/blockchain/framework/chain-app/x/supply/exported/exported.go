package exported

import (
	sdk "github.com/hashrs/blockchain/framework/chain-app/types"

	"github.com/hashrs/blockchain/framework/chain-app/x/auth/exported"
)

// ModuleAccountI defines an account interface for modules that hold tokens in an escrow
type ModuleAccountI interface {
	exported.Account

	GetName() string
	GetPermissions() []string
	HasPermission(string) bool
}

// SupplyI defines an inflationary supply interface for modules that handle
// token supply.
type SupplyI interface {
	GetTotal() sdk.Coins
	SetTotal(total sdk.Coins) SupplyI

	Inflate(amount sdk.Coins) SupplyI
	Deflate(amount sdk.Coins) SupplyI

	String() string
	ValidateBasic() error
}
