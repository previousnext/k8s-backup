package cronjob

import (
	"github.com/pkg/errors"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

// Deploy will create the CronJob and fallback to updating if it exists.
func Deploy(client *kubernetes.Clientset, cronjob *batchv2alpha1.CronJob) error {
	_, err := client.BatchV2alpha1().CronJobs(cronjob.ObjectMeta.Namespace).Create(cronjob)
	if kerrors.IsAlreadyExists(err) {
		_, err := client.BatchV2alpha1().CronJobs(cronjob.ObjectMeta.Namespace).Update(cronjob)
		if err != nil {
			return errors.Wrap(err, "failed to update")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	return nil
}
