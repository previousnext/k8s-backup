package config

import (
	"fmt"

	"github.com/pkg/errors"
)

// BucketURI for AWS S3.
func (p Config) BucketURI(namespace, strategy string) (string, error) {
	if p.Bucket == "" {
		return "", errors.New("not found: bucket")
	}

	return fmt.Sprintf("s3://%s/%s/%s", p.Bucket, namespace, strategy), nil
}
