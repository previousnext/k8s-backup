package pvc

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/previousnext/k8s-backup/config"
	skprcronjob "github.com/previousnext/skpr/utils/k8s/cronjob"
)

// Name of this backup strategy.
const Name = "pvc"

// Backoff after 1 failed attempt.
var Backoff int32 = 1

// Deploy backup strategies for PersistentVolumeClaims.
func Deploy(w io.Writer, client *kubernetes.Clientset, cfg config.Config) error {
	fmt.Fprintln(w, "Querying Kubernetes for PersistentVolumeClaims in namespace:", cfg.Namespace)

	pvcs, err := client.CoreV1().PersistentVolumeClaims(cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to lookup PersistentVolumeClaims")
	}

	for _, pvc := range pvcs.Items {
		fmt.Println("Syncing CronJob:", pvc.ObjectMeta.Namespace, "|", pvc.ObjectMeta.Name)

		cronjob, err := generateCronJob(pvc, cfg)
		if err != nil {
			return errors.Wrap(err, "failed to generate CronJob")
		}

		err = skprcronjob.Deploy(client, cronjob)
		if err != nil {
			return errors.Wrap(err, "failed to deploy CronJob")
		}
	}

	return nil
}

// Helper function to convert a PersistentVolumeClaim into a backup CronJob task.
func generateCronJob(pvc corev1.PersistentVolumeClaim, cfg config.Config) (*batchv1beta1.CronJob, error) {
	cronjob := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: pvc.ObjectMeta.Namespace,
			Name:      fmt.Sprintf("%s-%s", cfg.Prefix, pvc.ObjectMeta.Name),
		},
	}

	envvars, err := cfg.Credentials.EnvVars()
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to CronJob credentials")
	}

	resources, err := cfg.Resources.ResourceRequirements()
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to CronJob resources")
	}

	bucket, err := cfg.BucketURI(pvc.ObjectMeta.Namespace, Name)
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to CronJob bucket")
	}

	cronjob.Spec = batchv1beta1.CronJobSpec{
		Schedule:          cfg.Frequency,
		ConcurrencyPolicy: batchv1beta1.ForbidConcurrent,
		JobTemplate: batchv1beta1.JobTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: pvc.ObjectMeta.Namespace,
			},
			Spec: batchv1.JobSpec{
				BackoffLimit: &Backoff,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: pvc.ObjectMeta.Namespace,
					},
					Spec: corev1.PodSpec{
						RestartPolicy: "Never",
						Containers: []corev1.Container{
							{
								Name:  "sync",
								Image: cfg.Image,
								Command: []string{
									"aws",
									"s3",
									"sync",
									"/source/",
									fmt.Sprintf("%s/%s/", bucket, pvc.ObjectMeta.Name),
								},
								Env:       envvars,
								Resources: resources,
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
										ClaimName: pvc.ObjectMeta.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return cronjob, nil
}
