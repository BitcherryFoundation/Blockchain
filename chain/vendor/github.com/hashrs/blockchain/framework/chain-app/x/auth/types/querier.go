package types

import (
	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
)

// query endpoints supported by the auth Querier
const (
	QueryAccount = "account"
)

// QueryAccountParams defines the params for querying accounts.
type QueryAccountParams struct {
	Address sdk.AccAddress
}

// NewQueryAccountParams creates a new instance of QueryAccountParams.
func NewQueryAccountParams(addr sdk.AccAddress) QueryAccountParams {
	return QueryAccountParams{Address: addr}
}
