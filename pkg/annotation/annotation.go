package annotation

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Key is used for grouping backups in storage.
const Key = "k8s-backup.group"

// GetGroup will return the group based on annotations.
func GetGroup(meta metav1.ObjectMeta) (string, error) {
	if val, ok := meta.Annotations[Key]; ok {
		return val, nil
	}

	return "", fmt.Errorf("cannot find annotation: %s", Key)
}
