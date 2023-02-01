package config

import (
	"fmt"

	"github.com/argoproj/argo-rollouts/utils/defaults"
	"github.com/argoproj/argo-rollouts/utils/plugin/types"
	"github.com/ghodss/yaml"
	v1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	informers "k8s.io/client-go/informers/core/v1"
)

// Config is the in memory representation of the configmap with some additional fields/functions for ease of use.
type Config struct {
	configMap *v1.ConfigMap
	plugins   types.Plugin
}

var configMemoryCache *Config

// InitializeConfig initializes the in memory config and downloads the plugins to the filesystem. Subsequent calls to this function will return
// the same config object.
func InitializeConfig(configMapInformer informers.ConfigMapInformer, configMapName string, downloader FileDownloader) (*Config, error) {
	configMapCluster, err := configMapInformer.Lister().ConfigMaps(defaults.Namespace()).Get(configMapName)
	if err != nil {
		if k8errors.IsNotFound(err) {
			// If the configmap is not found, we return
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get configmap %s/%s: %w", defaults.Namespace(), configMapName, err)
	}

	plugins := types.Plugin{}
	if err = yaml.Unmarshal([]byte(configMapCluster.Data["plugins"]), &plugins); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugins while initializing: %w", err)
	}

	configMemoryCache = &Config{
		configMap: configMapCluster,
		plugins:   plugins,
	}

	if err := initMetricsPlugins(downloader); err != nil {
		return nil, err
	}
	return configMemoryCache, nil
}

// GetConfig returns the initialized in memory config object if it exists otherwise errors if InitializeConfig has not been called.
func GetConfig() (*Config, error) {
	if configMemoryCache == nil {
		return nil, fmt.Errorf("config not initialized, please initialize before use")
	}
	return configMemoryCache, nil
}

// GetMetricPluginsConfig returns the metric plugins configured in the configmap
func (c *Config) GetMetricPluginsConfig() []types.PluginItem {
	return configMemoryCache.plugins.Metrics
}
