package config

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ResourceRequirements Kubernetes object.
func (r Resources) ResourceRequirements() (corev1.ResourceRequirements, error) {
	cpu, err := resource.ParseQuantity(r.CPU)
	if err != nil {
		return corev1.ResourceRequirements{}, errors.Wrap(err, "failed to parse resource: cpu")
	}

	memory, err := resource.ParseQuantity(r.Memory)
	if err != nil {
		return corev1.ResourceRequirements{}, errors.Wrap(err, "failed to parse resource: memory")
	}

	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
	}, nil
}
