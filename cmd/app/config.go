package app

import (
	"encoding/json"
	"os"
)

// NodeInfo associated with core/database/nodes.go
type NodeInfo struct {
	Identity string `json:"identity"`

	// two slices of keys, base64-encoded, for JSON serialization
	SignKey    []string `json:"signKey"`
	EncryptKey []string `json:"encryptKey"`
	Address    []string `json:"address"`
}

func MarshalNodeInfoToJSON(node NodeInfo) ([]byte, error) {
	return json.Marshal(node)
}

func UnmarshalJSONToNodeInfo(b []byte) (NodeInfo, error) {
	var nodeInfo NodeInfo

	if err := json.Unmarshal(b, &nodeInfo); err != nil {
		return NodeInfo{}, err
	}

	return nodeInfo, nil
}

func ReadNodeInfo(path string) (NodeInfo, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return NodeInfo{}, err
	}

	return UnmarshalJSONToNodeInfo(b)
}

type Config struct {
	Protocols       []string `json:"Protocols"`
	AddressesIfTCP  []string `json:"addressesIfTCP"`
	LocalPortsIfTor []int    `json:"localPortsIfTor"`
	// todo other configs here
}

func MarshalBrConfigToJSON(cfg Config) ([]byte, error) {
	return json.Marshal(cfg)
}

func UnmarshalJSONToBrConfig(b []byte) (Config, error) {
	var BrConfig Config

	if err := json.Unmarshal(b, &BrConfig); err != nil {
		return Config{}, err
	}

	return BrConfig, nil
}

func ReadBrConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	return UnmarshalJSONToBrConfig(b)
}
