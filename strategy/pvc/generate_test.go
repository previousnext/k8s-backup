package pvc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/k8s-backup/config"
)

func TestGenerateCronJob(t *testing.T) {
	have, err := generateCronJob("dev", corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo",
			Name:      "bar",
		},
	}, config.Config{
		Namespace: "default",
		Image:     "previousnext/k8s-backup:latest",
		Prefix:    "test",
		CronSplit: 5,
		Bucket:    "test-bucket",
		Credentials: config.Credentials{
			ID:     "xxxxxxxxxxxxxxxx",
			Secret: "yyyyyyyyyyyyyyyy",
		},
		Resources: config.Resources{
			CPU:    "100m",
			Memory: "256Mi",
		},
	})
	assert.Nil(t, err)

	want := corev1.PodSpec{
		RestartPolicy: "Never",
		Containers: []corev1.Container{
			{
				Name:  "sync",
				Image: "previousnext/k8s-backup:latest",
				Command: []string{
					"aws",
					"s3",
					"sync",
					"/source/",
					"s3://test-bucket/foo/dev/pvc/bar/",
				},
				Env: []corev1.EnvVar{
					{
						Name:  "AWS_ACCESS_KEY_ID",
						Value: "xxxxxxxxxxxxxxxx",
					},
					{
						Name:  "AWS_SECRET_ACCESS_KEY",
						Value: "yyyyyyyyyyyyyyyy",
					},
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("256Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("256Mi"),
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "source",
						MountPath: "/source",
						ReadOnly:  true,
					},
				},
				ImagePullPolicy: "Always",
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: "source",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "bar",
					},
				},
			},
		},
	}

	assert.Equal(t, "foo", have.ObjectMeta.Namespace)
	assert.Equal(t, "test-pvc-bar", have.ObjectMeta.Name)
	assert.Equal(t, "* * * * *", have.Spec.Schedule)
	assert.Equal(t, want, have.Spec.JobTemplate.Spec.Template.Spec)
}
