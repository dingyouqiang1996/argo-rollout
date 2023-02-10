package plugin

import (
	"path/filepath"
	"testing"

	"github.com/argoproj/argo-rollouts/utils/defaults"

	"github.com/argoproj/argo-rollouts/utils/config"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetPluginLocation(t *testing.T) {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "argo-rollouts-config",
			Namespace: "argo-rollouts",
		},
		Data: map[string]string{"plugins": "metrics:\n  - plugin: argoproj-labs/http\n    pluginLocation: https://test/plugin\n  - plugin: argoproj-labs/http-sha\n    pluginLocation: https://test/plugin\n    pluginSha256: 74657374e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}
	client := fake.NewSimpleClientset(cm)

	_, err := config.InitializeConfig(client, "argo-rollouts-config")
	assert.NoError(t, err)

	location, err := GetPluginLocation("argoproj-labs/http")
	assert.NoError(t, err)
	fp, err := filepath.Abs(filepath.Join(defaults.DefaultRolloutPluginFolder, "argoproj-labs/http"))
	assert.NoError(t, err)
	assert.Equal(t, fp, location)
}

func TestGetPluginLocationNoNamedPlugin(t *testing.T) {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "argo-rollouts-config",
			Namespace: "argo-rollouts",
		},
		Data: map[string]string{"plugins": "metrics:\n  - plugin: argoproj-labs/http\n    pluginLocation: https://test/plugin\n  - plugin: argoproj-labs/http-sha\n    pluginLocation: https://test/plugin\n    pluginSha256: 74657374e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}
	client := fake.NewSimpleClientset(cm)

	_, err := config.InitializeConfig(client, "argo-rollouts-config")
	assert.NoError(t, err)

	location, err := GetPluginLocation("does-not-exist")
	assert.Error(t, err)
	assert.Equal(t, "plugin does-not-exist not configured in configmap", err.Error())
	assert.Equal(t, "", location)
}
