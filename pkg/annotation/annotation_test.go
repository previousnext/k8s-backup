package annotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetGroup(t *testing.T) {
	_, err := GetGroup(metav1.ObjectMeta{})
	assert.NotNil(t, err)

	annotation, err := GetGroup(metav1.ObjectMeta{
		Annotations: map[string]string{
			Key: "foo",
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "foo", annotation)
}
