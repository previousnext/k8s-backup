package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketURI(t *testing.T) {
	var cfg Config

	_, err := cfg.BucketURI("bar", "baz")
	assert.NotNil(t, err)

	// Need to set Bucket.
	cfg.Bucket = "foo"

	uri, err := cfg.BucketURI("bar", "baz")
	assert.Nil(t, err)
	assert.Equal(t, "s3://foo/bar/baz", uri)
}
