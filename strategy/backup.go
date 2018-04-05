package strategy

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	"github.com/previousnext/k8s-backup/config"
	"github.com/previousnext/k8s-backup/strategy/configmap/mysql"
	"github.com/previousnext/k8s-backup/strategy/pvc"
)

// Deploy backup strategies.
func Deploy(strategy []string, w io.Writer, client *kubernetes.Clientset, params config.Config) error {
	err := params.Validate()
	if err != nil {
		return errors.Wrap(err, "params are not valid")
	}

	for _, s := range strategy {
		if s == pvc.Name {
			err := pvc.Deploy(w, client, params)
			if err != nil {
				return errors.Wrap(err, "failed to sync PersistentVolumeClaim strategy")
			}
		}

		if s == mysql.Name {
			err := mysql.Deploy(w, client, params)
			if err != nil {
				return errors.Wrap(err, "failed to sync ConfigMap strategy")
			}
		}

		return fmt.Errorf("cannot find strategy: %s", s)
	}

	return nil
}
