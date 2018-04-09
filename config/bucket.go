package config

import (
	"fmt"

	"github.com/pkg/errors"
)

// BucketURI for AWS S3.
func (c Config) BucketURI(namespace, group string) (string, error) {
	if c.Bucket == "" {
		return "", errors.New("not found: bucket")
	}

	return fmt.Sprintf("s3://%s/%s/%s", c.Bucket, namespace, group), nil
}
