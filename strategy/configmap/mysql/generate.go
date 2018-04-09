package mysql

import (
	"fmt"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/k8s-backup/config"
)

// Name of this backup strategy.
const (
	KeyHostname = "mysql.hostname"
	KeyUsername = "mysql.username"
	KeyPassword = "mysql.password"
	KeyDatabase = "mysql.database"
)

// Helper function to extract a mysql connection from a ConfigMap key/value set.
func getMysqlConnection(configmap corev1.ConfigMap) (string, string, string, string, error) {
	if _, ok := configmap.Data[KeyHostname]; !ok {
		return "", "", "", "", fmt.Errorf("not found: %s", KeyHostname)
	}

	if _, ok := configmap.Data[KeyUsername]; !ok {
		return "", "", "", "", fmt.Errorf("not found: %s", KeyUsername)
	}

	if _, ok := configmap.Data[KeyPassword]; !ok {
		return "", "", "", "", fmt.Errorf("not found: %s", KeyPassword)
	}

	if _, ok := configmap.Data[KeyDatabase]; !ok {
		return "", "", "", "", fmt.Errorf("not found: %s", KeyDatabase)
	}

	return configmap.Data[KeyHostname], configmap.Data[KeyUsername], configmap.Data[KeyPassword], configmap.Data[KeyDatabase], nil
}

// Helper function to convert a PersistentVolumeClaim into a backup CronJob task.
func generateCronJob(group string, configmap corev1.ConfigMap, cfg config.Config) (*batchv1beta1.CronJob, error) {
	cronjob := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: configmap.ObjectMeta.Namespace,
			Name:      fmt.Sprintf("%s-configmap-mysql-%s", cfg.Prefix, configmap.ObjectMeta.Name),
		},
	}

	mysqlHost, mysqlUser, mysqlPass, mysqlName, err := getMysqlConnection(configmap)
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to get mysql connection from ConfigMap")
	}

	envs, err := cfg.Credentials.EnvVars()
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to CronJob credentials")
	}

	resources, err := cfg.Resources.ResourceRequirements()
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to CronJob resources")
	}

	bucket, err := cfg.BucketURI(configmap.ObjectMeta.Namespace, group)
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
				Namespace: configmap.ObjectMeta.Namespace,
			},
			Spec: batchv1.JobSpec{
				BackoffLimit: &backoff,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: configmap.ObjectMeta.Namespace,
					},
					Spec: corev1.PodSpec{
						RestartPolicy: "Never",
						InitContainers: []corev1.Container{
							{
								Name:  "dump",
								Image: cfg.Image,
								Command: []string{
									"/bin/sh", "-c",
								},
								Args: []string{
									fmt.Sprintf("mysqldump --host=%s --user=%s --pass=%s %s > /tmp/db.sql", mysqlHost, mysqlUser, mysqlPass, mysqlName),
								},
								Resources:       resources,
								ImagePullPolicy: "Always",
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "tmp",
										MountPath: "/tmp",
									},
								},
							},
						},
						Containers: []corev1.Container{
							{
								Name:  "push",
								Image: cfg.Image,
								Command: []string{
									"/bin/sh", "-c",
								},
								Args: []string{
									fmt.Sprintf("aws s3 cp /tmp/db.sql %s/configmap/mysql/%s.sql", bucket, mysqlName),
								},
								Env:             envs,
								Resources:       resources,
								ImagePullPolicy: "Always",
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "tmp",
										MountPath: "/tmp",
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "tmp",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{
										Medium: corev1.StorageMediumDefault,
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
