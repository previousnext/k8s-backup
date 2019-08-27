package mysql

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/previousnext/k8s-backup/config"
	"github.com/previousnext/k8s-backup/pkg/annotation"
	"github.com/previousnext/k8s-backup/pkg/cronutils"
	skprcronjob "github.com/previousnext/skpr/utils/k8s/cronjob"
)

// Name of this backup strategy.
const Name = "configmap_mysql"

// Deploy backup strategies for Mysql databases.
func Deploy(w io.Writer, client *kubernetes.Clientset, cfg config.Config) error {
	fmt.Fprintln(w, "Querying Kubernetes for ConfigMaps in namespace:", cfg.Namespace)

	configmaps, err := client.CoreV1().ConfigMaps(cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to lookup ConfigMaps")
	}

	schedule := cronutils.NewSplitter(cfg.CronSplit)

	for _, configmap := range configmaps.Items {
		fmt.Println("Syncing CronJob:", configmap.ObjectMeta.Namespace, "|", configmap.ObjectMeta.Name)

		group, err := annotation.GetGroup(configmap.ObjectMeta)
		if err != nil {
			fmt.Fprintln(w, "Skipping CronJob for ConfigMap:", err)
			continue
		}

		cronjob, err := generateCronJob(group, schedule.Increment(), configmap, cfg)
		if err != nil {
			fmt.Fprintln(w, "Skipping CronJob for ConfigMap:", err)
			continue
		}

		err = skprcronjob.Deploy(client, cronjob)
		if err != nil {
			return errors.Wrap(err, "failed to deploy CronJob")
		}
	}

	return nil
}
