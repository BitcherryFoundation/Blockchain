package lite

import (
	"github.com/hashrs/blockchain/core/consensus/dpos-pbft/types"
)

// Verifier checks the votes to make sure the block really is signed properly.
// Verifier must know the current or recent set of validitors by some other
// means.
type Verifier interface {
	Verify(sheader types.SignedHeader) error
	ChainID() string
}
