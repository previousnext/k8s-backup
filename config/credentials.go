package config

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

// EnvVars Kubernetes object.
func (c Credentials) EnvVars() ([]corev1.EnvVar, error) {
	if c.ID == "" {
		return []corev1.EnvVar{}, errors.New("not found: id")
	}

	if c.Secret == "" {
		return []corev1.EnvVar{}, errors.New("not found: secret")
	}

	return []corev1.EnvVar{
		{
			Name:  "AWS_ACCESS_KEY_ID",
			Value: c.ID,
		},
		{
			Name:  "AWS_SECRET_ACCESS_KEY",
			Value: c.Secret,
		},
	}, nil
}
