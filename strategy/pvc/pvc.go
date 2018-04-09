package pvc

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/previousnext/k8s-backup/config"
	"github.com/previousnext/k8s-backup/pkg/annotation"
	skprcronjob "github.com/previousnext/skpr/utils/k8s/cronjob"
)

// Name of this backup strategy.
const Name = "pvc"

// Deploy backup strategies for PersistentVolumeClaims.
func Deploy(w io.Writer, client *kubernetes.Clientset, cfg config.Config) error {
	fmt.Fprintln(w, "Querying Kubernetes for PersistentVolumeClaims in namespace:", cfg.Namespace)

	pvcs, err := client.CoreV1().PersistentVolumeClaims(cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to lookup PersistentVolumeClaims")
	}

	for _, pvc := range pvcs.Items {
		fmt.Println("Syncing CronJob:", pvc.ObjectMeta.Namespace, "|", pvc.ObjectMeta.Name)

		group, err := annotation.GetGroup(pvc.ObjectMeta)
		if err != nil {
			fmt.Fprintln(w, "Failed to get Group Annotation:", err)
			continue
		}

		cronjob, err := generateCronJob(group, pvc, cfg)
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
