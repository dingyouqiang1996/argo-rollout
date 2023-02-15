package client

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/argoproj/argo-rollouts/rollout/trafficrouting/plugin/rpc"
	"github.com/argoproj/argo-rollouts/utils/plugin"
	goPlugin "github.com/hashicorp/go-plugin"
)

type trafficPlugin struct {
	pluginClient map[string]*goPlugin.Client
	plugin       map[string]rpc.TrafficRouterPlugin
}

var pluginClients *trafficPlugin
var once sync.Once

// GetTrafficPlugin returns a singleton plugin client for the given traffic router plugin. Calling this multiple times
// returns the same plugin client instance for the plugin name defined in the rollout object.
func GetTrafficPlugin(pluginName string) (rpc.TrafficRouterPlugin, error) {
	once.Do(func() {
		pluginClients = &trafficPlugin{
			pluginClient: make(map[string]*goPlugin.Client),
			plugin:       make(map[string]rpc.TrafficRouterPlugin),
		}
	})
	plugin, err := pluginClients.startPlugin(pluginName)
	if err != nil {
		return nil, fmt.Errorf("unable to start plugin system: %w", err)
	}
	return plugin, nil
}

func (t *trafficPlugin) startPlugin(pluginName string) (rpc.TrafficRouterPlugin, error) {
	var handshakeConfig = goPlugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "ARGO_ROLLOUTS_RPC_PLUGIN",
		MagicCookieValue: "trafficrouter",
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]goPlugin.Plugin{
		"RpcTrafficRouterPlugin": &rpc.RpcTrafficRouterPlugin{},
	}

	if t.pluginClient[pluginName] == nil || t.pluginClient[pluginName].Exited() {
		pluginPath, err := plugin.GetPluginLocation(pluginName)
		if err != nil {
			return nil, fmt.Errorf("unable to find plugin (%s): %w", pluginName, err)
		}

		t.pluginClient[pluginName] = goPlugin.NewClient(&goPlugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         pluginMap,
			Cmd:             exec.Command(pluginPath),
			Managed:         true,
		})

		rpcClient, err := t.pluginClient[pluginName].Client()
		if err != nil {
			return nil, err
		}

		// Request the plugin
		plugin, err := rpcClient.Dispense("RpcTrafficRouterPlugin")
		if err != nil {
			return nil, err
		}
		t.plugin[pluginName] = plugin.(rpc.TrafficRouterPlugin)

		err = t.plugin[pluginName].InitPlugin()
		if err.Error() != "" {
			return nil, err
		}
	}

	client, err := t.pluginClient[pluginName].Client()
	if err != nil {
		return nil, err
	}
	if err := client.Ping(); err != nil {
		t.pluginClient[pluginName].Kill()
		t.pluginClient[pluginName] = nil
		return nil, fmt.Errorf("could not ping plugin will cleanup process so we can restart it next reconcile (%w)", err)
	}

	return t.plugin[pluginName], nil
}
