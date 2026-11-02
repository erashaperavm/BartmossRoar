package database

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
)

type ID = [32]byte

type NodeDB struct {
	database *DB
}

// NodeRecord construct your all neighbor nodes at a period
type NodeRecord struct {
	Identity ID

	// MldsaKey This will be more credible than Keys Table since neighbor-nodes-first principle
	// ptr type will be transformed to real data when writing into db
	// here impl slices since network delay, served for PoT,
	// because PoT lifetime is longer than a SignKey and
	// its update cycle is async
	SignKey    []crypto.PublicKey
	EncryptKey []crypto.SymmetricCipherKey

	// Multi-addresses for a node, which can be used for different network protocols or interfaces
	// or prevent a single point of failure, or for redundancy and load balancing
	Address []string
}

// RecordChannelNode stores a node record in the database.
// Important: The Identity field of the NodeRecord is used as the key for storage, and it is not stored in the value itself.
func (db *NodeDB) RecordChannelNode(node NodeRecord) error {
	nd := NodeRecord{
		SignKey:    node.SignKey,
		EncryptKey: node.EncryptKey,
		Address:    node.Address,
	}

	b, err := nd.Bytes()
	if err != nil {
		return fmt.Errorf("encode channel node: %w", err)
	}

	if err := db.database.Set(Nodes, node.Identity[:], b, 0); err != nil {
		return fmt.Errorf("store channel node: %w", err)
	}
	return nil
}

// GetChannelNode retrieves a node record from the database by its ID.
// Important: The returned NodeRecord does not include the Identity field, as it is already known from the ID used for retrieval.
func (db *NodeDB) GetChannelNode(id ID) (*NodeRecord, error) {
	b, err := db.database.Get(Nodes, id[:])
	if err != nil {
		return nil, fmt.Errorf("load channel node: %w", err)
	}

	node, err := FromGobBytesToNode(b)
	if err != nil {
		return nil, fmt.Errorf("decode channel node: %w", err)
	}

	return node, nil
}

func (node *NodeRecord) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(node)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrEncode, err)
	}

	return buf.Bytes(), nil
}

func FromGobBytesToNode(data []byte) (*NodeRecord, error) {
	var node NodeRecord
	buf := bytes.NewBuffer(data)

	dec := gob.NewDecoder(buf)
	err := dec.Decode(&node)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecode, err)
	}

	return &node, nil

}
