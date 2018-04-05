package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestEnvVars(t *testing.T) {
	var creds Credentials

	_, err := creds.EnvVars()
	assert.NotNil(t, err)

	// Need to set ID.
	creds.ID = "foo"

	_, err = creds.EnvVars()
	assert.NotNil(t, err)

	// Need to set Secret.
	creds.Secret = "bar"

	vars, err := creds.EnvVars()
	assert.Nil(t, err)
	assert.Equal(t, []corev1.EnvVar{
		{
			Name:  "AWS_ACCESS_KEY_ID",
			Value: "foo",
		},
		{
			Name:  "AWS_SECRET_ACCESS_KEY",
			Value: "bar",
		},
	}, vars)
}
