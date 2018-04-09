package config

import (
	"github.com/pkg/errors"
)

// Validate the config.
func (c Config) Validate() error {
	if c.Image == "" {
		return errors.New("not found: image")
	}

	if c.Prefix == "" {
		return errors.New("not found: prefix")
	}

	if c.Bucket == "" {
		return errors.New("not found: bucket")
	}

	if c.Frequency == "" {
		return errors.New("not found: frequency")
	}

	if c.Resources.CPU == "" {
		return errors.New("not found: cpu")
	}

	if c.Resources.Memory == "" {
		return errors.New("not found: resources: memory")
	}

	if c.Credentials.ID == "" {
		return errors.New("not found: credentials: id")
	}

	if c.Credentials.Secret == "" {
		return errors.New("not found: credentials: secret")
	}

	return nil
}
