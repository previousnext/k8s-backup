package pvc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/k8s-backup/config"
)

func TestGenerateCronJob(t *testing.T) {
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo",
			Name:      "bar",
		},
	}

	have, err := generateCronJob(pvc, config.Config{
		Namespace: "default",
		Image:     "previousnext/k8s-backup:latest",
		Prefix:    "test",
		Frequency: "* * * * *",
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

	want := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo",
			Name:      "test-bar",
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule:          "* * * * *",
			ConcurrencyPolicy: batchv1beta1.ForbidConcurrent,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo",
				},
				Spec: batchv1.JobSpec{
					BackoffLimit: &Backoff,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "foo",
						},
						Spec: corev1.PodSpec{
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
										"s3://test-bucket/foo/pvc/bar/",
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
						},
					},
				},
			},
		},
	}

	assert.Equal(t, want, have)
}