package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kingpin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/previousnext/k8s-backup/config"
	"github.com/previousnext/k8s-backup/strategy"
)

var (
	cliStrategies = kingpin.Flag("strategies", "Strategies to use for backing up state").Default("pvc,configmap_mysql").Envar("BACKUP_STRATEGIES").String()
	cliImage      = kingpin.Flag("image", "Image to use for backup strategies").Default("previousnext/k8s-backup:latest").Envar("BACKUP_IMAGE").String()
	cliNamespace  = kingpin.Flag("namespace", "Namespace to create backup CronJobs for PersistentVolumeClaims").Default(corev1.NamespaceAll).Envar("K8S_NAMESPACE").String()
	cliFrequency  = kingpin.Flag("frequency", "How often to run the CronJob").Default("@daily").Envar("BACKUP_FREQUENCY").String()
	cliPrefix     = kingpin.Flag("prefix", "Prefix to use for CronJob names").Default("k8s-backup").Envar("BACKUP_PREFIX").String()
	cliBucket     = kingpin.Flag("aws-bucket", "Bucket to sync PersistentVolumeClaims files").Default("k8s-backup").Envar("AWS_S3_BUCKET").String()
	cliCredID     = kingpin.Flag("aws-id", "Credentials to use when syncing PersistentVolumeClaim").Default("").Envar("AWS_ACCESS_KEY_ID").String()
	cliCredSecret = kingpin.Flag("aws-secret", "Credentials to use when syncing PersistentVolumeClaim").Default("").Envar("AWS_SECRET_ACCESS_KEY").String()
	cliCPU        = kingpin.Flag("cpu", "How much CPU to allocate for CronJobs").Default("100m").Envar("BACKUP_CPU").String()
	cliMemory     = kingpin.Flag("memory", "How much Memory to allocate for CronJobs").Default("256Mi").Envar("BACKUP_MEMORY").String()
)

func main() {
	kingpin.Parse()

	fmt.Println("Starting...")

	k8sconfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to Kubernetes")

	k8sclient, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting to deploy backup strategies")

	err = strategy.Deploy(strings.Split(*cliStrategies, ","), os.Stdout, k8sclient, config.Config{
		Image:     *cliImage,
		Namespace: *cliNamespace,
		Frequency: *cliFrequency,
		Prefix:    *cliPrefix,
		Bucket:    *cliBucket,
		Credentials: config.Credentials{
			ID:     *cliCredID,
			Secret: *cliCredSecret,
		},
		Resources: config.Resources{
			CPU:    *cliCPU,
			Memory: *cliMemory,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Finished deploying backup strategies")
}
