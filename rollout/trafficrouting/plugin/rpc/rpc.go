package rpc

import (
	"encoding/gob"
	"fmt"
	"net/rpc"

	"github.com/argoproj/argo-rollouts/utils/plugin/types"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/hashicorp/go-plugin"
)

type UpdateHashArgs struct {
	Rollout                v1alpha1.Rollout
	CanaryHash             string
	StableHash             string
	AdditionalDestinations []v1alpha1.WeightDestination
}

type SetWeightAndVerifyWeightArgs struct {
	Rollout                v1alpha1.Rollout
	DesiredWeight          int32
	AdditionalDestinations []v1alpha1.WeightDestination
}

type SetHeaderArgs struct {
	Rollout        v1alpha1.Rollout
	SetHeaderRoute v1alpha1.SetHeaderRoute
}

type SetMirrorArgs struct {
	Rollout        v1alpha1.Rollout
	SetMirrorRoute v1alpha1.SetMirrorRoute
}

type RemoveManagedRoutesArgs struct {
	Rollout v1alpha1.Rollout
}

type VerifyWeightResponse struct {
	Verified bool
	Err      types.RpcError
}

func init() {
	gob.RegisterName("UpdateHashArgs", new(UpdateHashArgs))
	gob.RegisterName("SetWeightAndVerifyWeightArgs", new(SetWeightAndVerifyWeightArgs))
	gob.RegisterName("SetHeaderArgs", new(SetHeaderArgs))
	gob.RegisterName("SetMirrorArgs", new(SetMirrorArgs))
	gob.RegisterName("RemoveManagedRoutesArgs", new(RemoveManagedRoutesArgs))
}

// TrafficRouterPlugin is the interface that we're exposing as a plugin. It needs to match metricproviders.Providers but we can
// not import that package because it would create a circular dependency.
type TrafficRouterPlugin interface {
	NewTrafficRouterPlugin() types.RpcError
	types.RpcTrafficRoutingReconciler
}

// TrafficRouterPluginRPC Here is an implementation that talks over RPC
type TrafficRouterPluginRPC struct{ client *rpc.Client }

// NewTrafficRouterPlugin this is the client aka the controller side function that calls the server side rpc (plugin)
// this gets called once during startup of the plugin and can be used to setup informers or k8s clients etc.
func (g *TrafficRouterPluginRPC) NewTrafficRouterPlugin() types.RpcError {
	var resp types.RpcError
	err := g.client.Call("Plugin.NewTrafficRouterPlugin", new(interface{}), &resp)
	if err != nil {
		return types.RpcError{ErrorString: err.Error()}
	}
	return resp
}

func (g *TrafficRouterPluginRPC) UpdateHash(rollout *v1alpha1.Rollout, canaryHash string, stableHash string, additionalDestinations []v1alpha1.WeightDestination) types.RpcError {
	var resp types.RpcError
	var args interface{} = UpdateHashArgs{
		Rollout:                *rollout,
		CanaryHash:             canaryHash,
		StableHash:             stableHash,
		AdditionalDestinations: additionalDestinations,
	}
	err := g.client.Call("Plugin.UpdateHash", &args, &resp)
	if err != nil {
		return types.RpcError{ErrorString: err.Error()}
	}
	return resp
}

func (g *TrafficRouterPluginRPC) SetWeight(rollout *v1alpha1.Rollout, desiredWeight int32, additionalDestinations []v1alpha1.WeightDestination) types.RpcError {
	var resp types.RpcError
	var args interface{} = SetWeightAndVerifyWeightArgs{
		Rollout:                *rollout,
		DesiredWeight:          desiredWeight,
		AdditionalDestinations: additionalDestinations,
	}
	err := g.client.Call("Plugin.SetWeight", &args, &resp)
	if err != nil {
		return types.RpcError{ErrorString: err.Error()}
	}
	return resp
}

func (g *TrafficRouterPluginRPC) SetHeaderRoute(rollout *v1alpha1.Rollout, setHeaderRoute *v1alpha1.SetHeaderRoute) types.RpcError {
	var resp types.RpcError
	var args interface{} = SetHeaderArgs{
		Rollout:        *rollout,
		SetHeaderRoute: *setHeaderRoute,
	}
	err := g.client.Call("Plugin.SetHeaderRoute", &args, &resp)
	if err != nil {
		return types.RpcError{ErrorString: err.Error()}
	}
	return resp
}

func (g *TrafficRouterPluginRPC) SetMirrorRoute(rollout *v1alpha1.Rollout, setMirrorRoute *v1alpha1.SetMirrorRoute) types.RpcError {
	var resp types.RpcError
	var args interface{} = SetMirrorArgs{
		Rollout:        *rollout,
		SetMirrorRoute: *setMirrorRoute,
	}
	err := g.client.Call("Plugin.SetMirrorRoute", &args, &resp)
	if err != nil {
		return types.RpcError{ErrorString: err.Error()}
	}
	return resp
}

func (g *TrafficRouterPluginRPC) Type() string {
	var resp string
	err := g.client.Call("Plugin.Type", new(interface{}), &resp)
	if err != nil {
		return err.Error()
	}

	return resp
}

func (g *TrafficRouterPluginRPC) VerifyWeight(rollout *v1alpha1.Rollout, desiredWeight int32, additionalDestinations []v1alpha1.WeightDestination) (*bool, types.RpcError) {
	var resp VerifyWeightResponse
	var args interface{} = SetWeightAndVerifyWeightArgs{
		Rollout:                *rollout,
		DesiredWeight:          desiredWeight,
		AdditionalDestinations: additionalDestinations,
	}
	err := g.client.Call("Plugin.VerifyWeight", &args, &resp)
	if err != nil {
		return nil, types.RpcError{ErrorString: err.Error()}
	}
	return &resp.Verified, resp.Err
}

func (g *TrafficRouterPluginRPC) RemoveManagedRoutes(rollout *v1alpha1.Rollout) types.RpcError {
	var resp types.RpcError
	var args interface{} = RemoveManagedRoutesArgs{
		Rollout: *rollout,
	}
	err := g.client.Call("Plugin.RemoveManagedRoutes", &args, &resp)
	if err != nil {
		return types.RpcError{ErrorString: err.Error()}
	}
	return resp
}

// TrafficRouterRPCServer Here is the RPC server that MetricsPluginRPC talks to, conforming to
// the requirements of net/rpc
type TrafficRouterRPCServer struct {
	// This is the real implementation
	Impl TrafficRouterPlugin
}

func (s *TrafficRouterRPCServer) NewTrafficRouterPlugin(args interface{}, resp *types.RpcError) error {
	*resp = s.Impl.NewTrafficRouterPlugin()
	return nil
}

func (s *TrafficRouterRPCServer) UpdateHash(args interface{}, resp *types.RpcError) error {
	runArgs, ok := args.(*UpdateHashArgs)
	if !ok {
		return fmt.Errorf("invalid args %s", args)
	}
	*resp = s.Impl.UpdateHash(&runArgs.Rollout, runArgs.CanaryHash, runArgs.StableHash, runArgs.AdditionalDestinations)
	return nil
}

func (s *TrafficRouterRPCServer) SetWeight(args interface{}, resp *types.RpcError) error {
	setWeigthArgs, ok := args.(*SetWeightAndVerifyWeightArgs)
	if !ok {
		return fmt.Errorf("invalid args %s", args)
	}
	*resp = s.Impl.SetWeight(&setWeigthArgs.Rollout, setWeigthArgs.DesiredWeight, setWeigthArgs.AdditionalDestinations)
	return nil
}

func (s *TrafficRouterRPCServer) SetHeaderRoute(args interface{}, resp *types.RpcError) error {
	setHeaderArgs, ok := args.(*SetHeaderArgs)
	if !ok {
		return fmt.Errorf("invalid args %s", args)
	}
	*resp = s.Impl.SetHeaderRoute(&setHeaderArgs.Rollout, &setHeaderArgs.SetHeaderRoute)
	return nil
}

func (s *TrafficRouterRPCServer) SetMirrorRoute(args interface{}, resp *types.RpcError) error {
	setMirrorArgs, ok := args.(*SetMirrorArgs)
	if !ok {
		return fmt.Errorf("invalid args %s", args)
	}
	*resp = s.Impl.SetMirrorRoute(&setMirrorArgs.Rollout, &setMirrorArgs.SetMirrorRoute)
	return nil
}

func (s *TrafficRouterRPCServer) Type(args interface{}, resp *string) error {
	*resp = s.Impl.Type()
	return nil
}

func (s *TrafficRouterRPCServer) VerifyWeight(args interface{}, resp *VerifyWeightResponse) error {
	verifyWeightArgs, ok := args.(*SetWeightAndVerifyWeightArgs)
	if !ok {
		return fmt.Errorf("invalid args %s", args)
	}
	verified, err := s.Impl.VerifyWeight(&verifyWeightArgs.Rollout, verifyWeightArgs.DesiredWeight, verifyWeightArgs.AdditionalDestinations)
	*resp = VerifyWeightResponse{
		Verified: *verified,
		Err:      err,
	}
	return nil
}

func (s *TrafficRouterRPCServer) RemoveManagedRoutes(args interface{}, resp *types.RpcError) error {
	removeManagedRoutesArgs, ok := args.(*RemoveManagedRoutesArgs)
	if !ok {
		return fmt.Errorf("invalid args %s", args)
	}
	*resp = s.Impl.RemoveManagedRoutes(&removeManagedRoutesArgs.Rollout)
	return nil
}

// RpcTrafficRouterPlugin This is the implementation of plugin.Plugin so we can serve/consume
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a MetricsRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return MetricsPluginRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type RpcTrafficRouterPlugin struct {
	// Impl Injection
	Impl TrafficRouterPlugin
}

func (p *RpcTrafficRouterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &TrafficRouterRPCServer{Impl: p.Impl}, nil
}

func (RpcTrafficRouterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &TrafficRouterPluginRPC{client: c}, nil
}
