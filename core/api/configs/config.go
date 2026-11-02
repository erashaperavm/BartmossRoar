package configs

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ApiConfig struct {
	// todo, when needs
}

func MarshalApiConfigToYAML(ac ApiConfig) ([]byte, error) {
	return yaml.Marshal(ac)
}

func UnmarshalYAMLToApiConfig(b []byte) (ApiConfig, error) {
	var ac ApiConfig
	err := yaml.Unmarshal(b, ac)
	if err != nil {
		return ApiConfig{}, err
	}

	return ac, nil
}

func ReadApiConfig(path string) (ApiConfig, error) {
	b, err := os.ReadFile(filepath.Dir(path))
	if err != nil {
		return ApiConfig{}, err
	}

	return UnmarshalYAMLToApiConfig(b)
}
