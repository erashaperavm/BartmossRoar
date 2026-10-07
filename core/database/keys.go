package database

import "github.com/cloudflare/circl/sign/mldsa/mldsa44"

type ID [32]byte

// CommonNode construct the all you know nodes db
type CommonNode struct {
	identity ID

	// MaybePubKey is not always correct since the network delay and consensus
	// ptr type will be transformed to real data when writing into db
	maybePubKey *mldsa44.PublicKey
}
