package fixtures

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	rov1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	unstructuredutil "github.com/argoproj/argo-rollouts/utils/unstructured"
)

type Given struct {
	Common
}

// RolloutObjects sets up the rollout objects for the test environment given a YAML string or file path:
// 1. A file name if it starts with "@"
// 2. Raw YAML.
func (g *Given) RolloutObjects(text string) *Given {
	g.t.Helper()
	yamlBytes := g.yamlBytes(text)

	// Some E2E AnalysisTemplates use http://kubernetes.default.svc/version as a fake metric provider.
	// This doesn't work outside the cluster, so the following replaces it with the host from the
	// rest config.
	yamlString := strings.ReplaceAll(string(yamlBytes), "https://kubernetes.default.svc", g.kubernetesHost)

	objs, err := unstructuredutil.SplitYAML(yamlString)
	g.CheckError(err)
	for _, obj := range objs {
		labels := obj.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		if E2ELabelValueInstanceID != "" {
			labels[rov1.LabelKeyControllerInstanceID] = E2ELabelValueInstanceID
		}
		labels[E2ELabelKeyTestName] = g.t.Name()
		obj.SetLabels(labels)

		if obj.GetKind() == "Rollout" {
			if g.rollout != nil {
				g.t.Fatal("multiple rollouts specified")
			}
			g.rollout = &rov1.Rollout{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &g.rollout)
			g.CheckError(err)
			g.log = g.log.WithField("rollout", g.rollout.Name)

			if E2EPodDelay > 0 {
				// Add postStart/preStop handlers to slow down readiness/termination
				sleepHandler := corev1.Handler{
					Exec: &corev1.ExecAction{
						Command: []string{"sleep", strconv.Itoa(E2EPodDelay)},
					},
				}
				g.rollout.Spec.Template.Spec.Containers[0].Lifecycle = &corev1.Lifecycle{
					PostStart: &sleepHandler,
					PreStop:   &sleepHandler,
				}
			}
		} else {
			// other non-rollout objects
			g.objects = append(g.objects, obj)
		}
	}
	return g
}

func (g *Given) RolloutTemplate(text, name string) *Given {
	yamlBytes := g.yamlBytes(text)
	newText := strings.ReplaceAll(string(yamlBytes), "REPLACEME", name)
	return g.RolloutObjects(newText)
}

func (g *Given) yamlBytes(text string) []byte {
	var yamlBytes []byte
	var err error
	if strings.HasPrefix(text, "@") {
		file := strings.TrimPrefix(text, "@")
		yamlBytes, err = ioutil.ReadFile(file)
		g.CheckError(err)
	} else {
		yamlBytes = []byte(text)
	}
	return yamlBytes
}

func (g *Given) SetSteps(text string) *Given {
	steps := make([]rov1.CanaryStep, 0)
	err := yaml.Unmarshal([]byte(text), &steps)
	g.CheckError(err)
	g.rollout.Spec.Strategy.Canary.Steps = steps
	return g
}

// HealthyRollout is a convenience around creating a rollout and waiting for it to become healthy
func (g *Given) HealthyRollout(text string) *Given {
	return g.RolloutObjects(text).
		When().
		ApplyManifests().
		WaitForRolloutStatus("Healthy").
		Given()
}

func (g *Given) When() *When {
	return &When{
		Common: g.Common,
	}
}
