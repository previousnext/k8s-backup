package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestResourceRequirements(t *testing.T) {
	var resources Resources

	_, err := resources.ResourceRequirements()
	assert.NotNil(t, err)

	// Need to set CPU.
	resources.CPU = "100m"

	_, err = resources.ResourceRequirements()
	assert.NotNil(t, err)

	// Need to set Memory.
	resources.Memory = "256Mi"

	req, err := resources.ResourceRequirements()
	assert.Nil(t, err)
	assert.Equal(t, corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("256Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("256Mi"),
		},
	}, req)
}
