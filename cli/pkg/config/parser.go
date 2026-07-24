package config

import (
	"fmt"
	"io"
	"gopkg.in/yaml.v3"
)

type Policy struct {
	Version int `yaml:"version"`
	Config  struct {
		EnforcementMode string `yaml:"enforcement_mode"`
		Network         struct {
			Allow []string `yaml:"allow"`
			Deny  []string `yaml:"deny"`
		} `yaml:"network"`
	} `yaml:"policy"`
}

// ParsePolicy securely unmarshals a StuntDouble yaml policy and validates its structure
func ParsePolicy(r io.Reader) (*Policy, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var p Policy
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	if p.Version == 0 {
		return nil, fmt.Errorf("invalid policy: missing or invalid version")
	}

	return &p, nil
}
