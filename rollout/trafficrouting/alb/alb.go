package alb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	extensionslisters "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/utils/pointer"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/argoproj/argo-rollouts/rollout/trafficrouting"
	"github.com/argoproj/argo-rollouts/utils/aws"
	"github.com/argoproj/argo-rollouts/utils/diff"
	ingressutil "github.com/argoproj/argo-rollouts/utils/ingress"
	jsonutil "github.com/argoproj/argo-rollouts/utils/json"
	logutil "github.com/argoproj/argo-rollouts/utils/log"
	"github.com/argoproj/argo-rollouts/utils/record"
)

const (
	// Type holds this controller type
	Type = "ALB"
)

// ReconcilerConfig describes static configuration data for the ALB Ingress reconciler
type ReconcilerConfig struct {
	Rollout        *v1alpha1.Rollout
	Client         kubernetes.Interface
	Recorder       record.EventRecorder
	ControllerKind schema.GroupVersionKind
	IngressLister  extensionslisters.IngressLister
	VerifyWeight   *bool
}

// Reconciler holds required fields to reconcile ALB Ingress resources
type Reconciler struct {
	cfg ReconcilerConfig
	log *logrus.Entry
	aws aws.Client
}

var (
	defaultVerifyWeight = false
)

// SetDefaultVerifyWeight sets the default setWeight verification when instantiating the reconciler
func SetDefaultVerifyWeight(b bool) {
	defaultVerifyWeight = b
}

// NewReconciler returns a reconciler struct that brings the ALB Ingress into the desired state
func NewReconciler(cfg ReconcilerConfig) (*Reconciler, error) {
	awsClient, err := aws.NewClient()
	if err != nil {
		return nil, err
	}
	reconciler := Reconciler{
		cfg: cfg,
		log: logutil.WithRollout(cfg.Rollout).WithField(logutil.IngressKey, cfg.Rollout.Spec.Strategy.Canary.TrafficRouting.ALB.Ingress),
		aws: awsClient,
	}
	return &reconciler, nil
}

// Type indicates this reconciler is an ALB ingress reconciler
func (r *Reconciler) Type() string {
	return Type
}

// SetWeight modifies ALB Ingress resources to reach desired state
func (r *Reconciler) SetWeight(desiredWeight int32, additionalDestinations ...trafficrouting.WeightDestination) error {
	ctx := context.TODO()
	rollout := r.cfg.Rollout
	ingressName := rollout.Spec.Strategy.Canary.TrafficRouting.ALB.Ingress
	ingress, err := r.cfg.IngressLister.Ingresses(rollout.Namespace).Get(ingressName)
	if err != nil {
		return err
	}
	actionService := r.cfg.Rollout.Spec.Strategy.Canary.StableService
	if r.cfg.Rollout.Spec.Strategy.Canary.TrafficRouting.ALB.RootService != "" {
		actionService = r.cfg.Rollout.Spec.Strategy.Canary.TrafficRouting.ALB.RootService
	}
	port := r.cfg.Rollout.Spec.Strategy.Canary.TrafficRouting.ALB.ServicePort
	if !ingressutil.HasRuleWithService(ingress, actionService) {
		return fmt.Errorf("ingress does not have service `%s` in rules", actionService)
	}

	desired, err := getDesiredAnnotations(ingress, rollout, port, desiredWeight, additionalDestinations...)
	if err != nil {
		return err
	}
	patch, modified, err := calculatePatch(ingress, desired)
	if err != nil {
		return nil
	}
	if !modified {
		r.log.Info("no changes to the ALB Ingress")
		return nil
	}
	r.log.WithField("patch", string(patch)).Debug("applying ALB Ingress patch")
	r.log.WithField("desiredWeight", desiredWeight).Info("updating ALB Ingress")
	r.cfg.Recorder.Eventf(r.cfg.Rollout, record.EventOptions{EventReason: "PatchingALBIngress"}, "Updating Ingress `%s` to desiredWeight '%d'", ingressName, desiredWeight)
	_, err = r.cfg.Client.ExtensionsV1beta1().Ingresses(ingress.Namespace).Patch(ctx, ingress.Name, types.MergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		r.log.WithField("err", err.Error()).Error("error patching alb ingress")
		return fmt.Errorf("error patching alb ingress `%s`: %v", ingressName, err)
	}
	return nil
}

func (r *Reconciler) shouldVerifyWeight() bool {
	if r.cfg.VerifyWeight != nil {
		return *r.cfg.VerifyWeight
	}
	return defaultVerifyWeight
}

func (r *Reconciler) VerifyWeight(desiredWeight int32, additionalDestinations ...trafficrouting.WeightDestination) (bool, error) {
	if !r.shouldVerifyWeight() {
		return true, nil
	}
	ctx := context.TODO()
	rollout := r.cfg.Rollout
	ingressName := rollout.Spec.Strategy.Canary.TrafficRouting.ALB.Ingress
	ingress, err := r.cfg.IngressLister.Ingresses(rollout.Namespace).Get(ingressName)
	if err != nil {
		return false, err
	}
	resourceIDToService := map[string]string{}

	canaryService := rollout.Spec.Strategy.Canary.CanaryService
	canaryResourceID := aws.BuildV2TargetGroupID(rollout.Namespace, ingress.Name, canaryService, rollout.Spec.Strategy.Canary.TrafficRouting.ALB.ServicePort)
	resourceIDToService[canaryResourceID] = canaryService

	for _, dest := range additionalDestinations {
		resourceID := aws.BuildV2TargetGroupID(rollout.Namespace, ingress.Name, dest.ServiceName, rollout.Spec.Strategy.Canary.TrafficRouting.ALB.ServicePort)
		resourceIDToService[resourceID] = dest.ServiceName
	}

	if len(ingress.Status.LoadBalancer.Ingress) == 0 {
		r.log.Infof("LoadBalancer not yet allocated")
	}
	for _, lbIngress := range ingress.Status.LoadBalancer.Ingress {
		if lbIngress.Hostname == "" {
			continue
		}
		lb, err := r.aws.FindLoadBalancerByDNSName(ctx, lbIngress.Hostname)
		if err != nil {
			return false, err
		}
		if lb == nil || lb.LoadBalancerArn == nil {
			r.log.Infof("LoadBalancer %s not found", lbIngress.Hostname)
			return false, nil
		}
		lbTargetGroups, err := r.aws.GetTargetGroupMetadata(ctx, *lb.LoadBalancerArn)
		if err != nil {
			return false, err
		}
		logCtx := r.log.WithField("lb", *lb.LoadBalancerArn)
		for _, tg := range lbTargetGroups {
			//if tg.Tags[aws.AWSLoadBalancerV2TagKeyResourceID] == resourceID {
			resourceID := tg.Tags[aws.AWSLoadBalancerV2TagKeyResourceID]
			if svcName, ok := resourceIDToService[resourceID]; ok {
				if tg.Weight != nil {
					logCtx := logCtx.WithField("tg", *tg.TargetGroupArn)
					if resourceID == canaryResourceID {
						logCtx.Infof("canary weight of %s (desired: %d, current: %d)", resourceID, desiredWeight, *tg.Weight)
					} else {
						logCtx.Infof("%s weight of %s (desired: %d, current: %d)", svcName, resourceID, desiredWeight, *tg.Weight)
					}
					return *tg.Weight == desiredWeight, nil
				}
			}
		}
	}
	return false, nil
}

func calculatePatch(current *extensionsv1beta1.Ingress, desiredAnnotations map[string]string) ([]byte, bool, error) {
	// only compare Annotations
	return diff.CreateTwoWayMergePatch(
		&extensionsv1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: current.Annotations,
			},
		},
		&extensionsv1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: desiredAnnotations,
			},
		}, extensionsv1beta1.Ingress{})
}

func getForwardActionString(r *v1alpha1.Rollout, port int32, desiredWeight int32, additionalDestinations ...trafficrouting.WeightDestination) string {
	stableService := r.Spec.Strategy.Canary.StableService
	canaryService := r.Spec.Strategy.Canary.CanaryService
	portStr := strconv.Itoa(int(port))
	stableWeight := int32(100)
	targetGroups := make([]ingressutil.ALBTargetGroup, 0)
	// create target group for canary
	targetGroups = append(targetGroups, ingressutil.ALBTargetGroup{
		ServiceName: canaryService,
		ServicePort: portStr,
		Weight:      pointer.Int64Ptr(int64(desiredWeight)),
	})
	// update stableWeight
	stableWeight -= desiredWeight

	for _, dest := range additionalDestinations {
		// Create target group for each additional destination
		targetGroups = append(targetGroups, ingressutil.ALBTargetGroup{
			ServiceName: dest.ServiceName,
			ServicePort: portStr,
			Weight:      pointer.Int64Ptr(int64(dest.Weight)),
		})
		stableWeight -= dest.Weight
	}

	// Create target group for stable with updated stableWeight
	targetGroups = append(targetGroups, ingressutil.ALBTargetGroup{
		ServiceName: stableService,
		ServicePort: portStr,
		Weight:      pointer.Int64Ptr(int64(stableWeight)),
	})

	action := ingressutil.ALBAction{
		Type: "forward",
		ForwardConfig: ingressutil.ALBForwardConfig{
			TargetGroups: targetGroups,
		},
	}
	bytes := jsonutil.MustMarshal(action)
	return string(bytes)
}

func getDesiredAnnotations(current *extensionsv1beta1.Ingress, r *v1alpha1.Rollout, port int32, desiredWeight int32, additionalDestinations ...trafficrouting.WeightDestination) (map[string]string, error) {
	desired := current.DeepCopy().Annotations
	key := ingressutil.ALBActionAnnotationKey(r)
	desired[key] = getForwardActionString(r, port, desiredWeight, additionalDestinations...)
	m, err := ingressutil.NewManagedALBActions(desired[ingressutil.ManagedActionsAnnotation])
	if err != nil {
		return nil, err
	}
	m[r.Name] = key
	desired[ingressutil.ManagedActionsAnnotation] = m.String()
	return desired, nil
}

// UpdateHash informs a traffic routing reconciler about new canary/stable pod hashes
func (r *Reconciler) UpdateHash(canaryHash, stableHash string) error {
	return nil
}
