package pvc

import (
	"fmt"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/k8s-backup/config"
)

// Helper function to convert a PersistentVolumeClaim into a backup CronJob task.
func generateCronJob(group string, pvc corev1.PersistentVolumeClaim, cfg config.Config) (*batchv1beta1.CronJob, error) {
	cronjob := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: pvc.ObjectMeta.Namespace,
			Name:      fmt.Sprintf("%s-pvc-%s", cfg.Prefix, pvc.ObjectMeta.Name),
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

	bucket, err := cfg.BucketURI(pvc.ObjectMeta.Namespace, group)
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to CronJob bucket")
	}

	var (
		// Backoff determines how many times the build fails before it does not get recreated.
		// This is set to 2 for:
		//  * Generally a fail will be OOMKiller not happy with how much memory awscli is using,
		//    we shouldn't run builds over and over again, they will keep failing.
		//  * This amount allows for any "transient" issues that could be fixed with a rerun.
		backoff int32 = 2
		// CronJobs will have to start within 30 min.
		// https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#starting-deadline-seconds
		deadline int64 = 1800
	)

	cronjob.Spec = batchv1beta1.CronJobSpec{
		Schedule:                cfg.Frequency,
		ConcurrencyPolicy:       batchv1beta1.ForbidConcurrent,
		StartingDeadlineSeconds: &deadline,
		JobTemplate: batchv1beta1.JobTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: pvc.ObjectMeta.Namespace,
			},
			Spec: batchv1.JobSpec{
				BackoffLimit: &backoff,
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
									fmt.Sprintf("%s/pvc/%s/", bucket, pvc.ObjectMeta.Name),
								},
								Env:       envvars,
								Resources: resources,
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
