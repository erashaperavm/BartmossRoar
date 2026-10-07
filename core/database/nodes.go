package database

import "github.com/erashaperavm/BartmossRoar/core/crypto"

// ChannelNode construct your all neighbor nodes at a period
type ChannelNode struct {
	identity ID

	// MldsaKey This will be more credible than Keys Table since neighbor-nodes-first principle
	// ptr type will be transformed to real data when writing into db
	signKey    *crypto.PublicKey
	encryptKey *crypto.Key
}
