package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	var cfg Config

	err := cfg.Validate()
	assert.NotNil(t, err)

	// Need to set Image.
	cfg.Image = "previousnext/k8s-backup:latest"

	err = cfg.Validate()
	assert.NotNil(t, err)

	// Need to set Prefix.
	cfg.Prefix = "k8s-backup"

	err = cfg.Validate()
	assert.NotNil(t, err)

	// Need to set Bucket./
	cfg.Bucket = "k8s-backup"

	err = cfg.Validate()
	assert.NotNil(t, err)

	// Need to set Frequency.
	cfg.Frequency = "* * * * *"

	err = cfg.Validate()
	assert.NotNil(t, err)

	// Need to set Resources.
	cfg.Resources = Resources{
		CPU:    "100m",
		Memory: "256Mi",
	}

	err = cfg.Validate()
	assert.NotNil(t, err)

	// Need to set Credentials.
	cfg.Credentials = Credentials{
		ID:     "xxxxxxxxxxxxx",
		Secret: "yyyyyyyyyyyyy",
	}

	err = cfg.Validate()
	assert.Nil(t, err)
}
