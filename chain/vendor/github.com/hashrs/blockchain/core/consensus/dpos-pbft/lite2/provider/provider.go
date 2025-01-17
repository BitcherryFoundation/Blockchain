package provider

import "github.com/hashrs/blockchain/core/consensus/dpos-pbft/types"

// Provider provides information for the lite client to sync (verification
// happens in the client).
type Provider interface {
	// ChainID returns the blockchain ID.
	ChainID() string

	// SignedHeader returns the SignedHeader that corresponds to the given
	// height.
	//
	// 0 - the latest.
	// height must be >= 0.
	//
	// If the provider fails to fetch the SignedHeader due to the IO or other
	// issues, an error will be returned.
	// If there's no SignedHeader for the given height, ErrSignedHeaderNotFound
	// will be returned.
	SignedHeader(height int64) (*types.SignedHeader, error)

	// ValidatorSet returns the ValidatorSet that corresponds to height.
	//
	// 0 - the latest.
	// height must be >= 0.
	//
	// If the provider fails to fetch the ValidatorSet due to the IO or other
	// issues, an error will be returned.
	// If there's no ValidatorSet for the given height, ErrValidatorSetNotFound
	// will be returned.
	ValidatorSet(height int64) (*types.ValidatorSet, error)
}
