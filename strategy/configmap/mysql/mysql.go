package mysql

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
const (
	Name        = "configmap_mysql"
	KeyHostname = "mysql.host"
	KeyUsername = "mysql.user"
	KeyPassword = "mysql.pass"
	KeyDatabase = "mysql.database"
)

// Backoff after 1 failed attempt.
var Backoff int32 = 1

// Deploy backup strategies for Mysql databases.
func Deploy(w io.Writer, client *kubernetes.Clientset, cfg config.Config) error {
	fmt.Fprintln(w, "Querying Kubernetes for ConfigMaps in namespace:", cfg.Namespace)

	configmaps, err := client.CoreV1().ConfigMaps(cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to lookup ConfigMaps")
	}

	for _, configmap := range configmaps.Items {
		fmt.Println("Syncing CronJob:", configmap.ObjectMeta.Namespace, "|", configmap.ObjectMeta.Name)

		// @todo, Filter by values.
		host, user, pass, db, err := extractMysqlConnection(configmap)
		if err != nil {
			fmt.Fprintln(w, "Skipping CronJob for ConfigMap:", configmap.ObjectMeta.Namespace, "|", configmap.ObjectMeta.Name)
			continue
		}

		cronjob, err := generateCronJob(configmap.ObjectMeta.Namespace, configmap.ObjectMeta.Name, host, user, pass, db, cfg)
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

func extractMysqlConnection(configmap corev1.ConfigMap) (string, string, string, string, error) {
	var (
		hostname string
		username string
		password string
		database string
	)

	if val, ok := configmap.Data[KeyHostname]; ok {
		hostname = val
	}

	if val, ok := configmap.Data[KeyUsername]; ok {
		username = val
	}

	if val, ok := configmap.Data[KeyPassword]; ok {
		password = val
	}

	if val, ok := configmap.Data[KeyDatabase]; ok {
		database = val
	}

	return hostname, username, password, database, nil
}

// Helper function to convert a PersistentVolumeClaim into a backup CronJob task.
func generateCronJob(namespace, name, hostname, username, password, database string, cfg config.Config) (*batchv1beta1.CronJob, error) {
	cronjob := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s-mysql-db-%s-%s", cfg.Prefix, name, database),
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

	bucket, err := cfg.BucketURI(namespace, Name)
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to CronJob bucket")
	}

	cronjob.Spec = batchv1beta1.CronJobSpec{
		Schedule:          cfg.Frequency,
		ConcurrencyPolicy: batchv1beta1.ForbidConcurrent,
		JobTemplate: batchv1beta1.JobTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
			},
			Spec: batchv1.JobSpec{
				BackoffLimit: &Backoff,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
					},
					Spec: corev1.PodSpec{
						RestartPolicy: "Never",
						Containers: []corev1.Container{
							{
								Name:  "sync",
								Image: cfg.Image,
								Command: []string{
									"mysqldump",
									"--host", hostname,
									"--user", username,
									"--pass", password,
									"--databases", database,
									"--result-file", "/tmp/db.sql",
									"&&",
									"aws", "s3", "sync", "/tmp/db.sql", fmt.Sprintf("%s/%s/%s.sql", bucket, name, database),
								},
								Env:             envvars,
								Resources:       resources,
								ImagePullPolicy: "Always",
							},
						},
					},
				},
			},
		},
	}

	return cronjob, nil
}
