package restart

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	clientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned/typed/rollouts/v1alpha1"
	"github.com/argoproj/argo-rollouts/pkg/kubectl-argo-rollouts/options"
)

const (
	restartExample = `
	# Restart the pods of a rollout in now
	%[1]s restart ROLLOUT_NAME

	# Restart the pods of a rollout in ten seconds
	%[1]s restart ROLLOUT_NAME --in 10s`

	restartPatch = `{
	"spec": {
		"restartAt": "%s"
	}
}`
)

func NewCmdRestart(o *options.ArgoRolloutsOptions) *cobra.Command {
	var (
		in string
	)
	var cmd = &cobra.Command{
		Use:          "restart ROLLOUT",
		Short:        "Restart the pods of a rollout",
		Example:      o.Example(restartExample),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return o.UsageErr(c)
			}
			restartAt := o.Now().UTC()
			if in != "" {
				duration, err := v1alpha1.DurationString(in).Duration()
				if err != nil {
					panic(err)
				}
				restartAt = restartAt.Add(duration)
			} else {
				in = "0s"
			}
			name := args[0]
			rolloutIf := o.RolloutsClientset().ArgoprojV1alpha1().Rollouts(o.Namespace())
			ro, err := RestartRollout(rolloutIf, name, &restartAt)
			if err != nil {
				return err
			}
			fmt.Fprintf(o.Out, "rollout '%s' restarts in %s\n", ro.Name, in)
			return nil
		},
	}
	cmd.Flags().StringVarP(&in, "in", "i", "", "Amount of time before a restart. (e.g. 30s, 5m, 1h)")
	return cmd
}

// RestartRollout restarts a rollout
func RestartRollout(rolloutIf clientset.RolloutInterface, name string, restartAt *time.Time) (*v1alpha1.Rollout, error) {
	if restartAt == nil {
		t := time.Now().UTC()
		restartAt = &t
	}
	patch := fmt.Sprintf(restartPatch, restartAt.Format(time.RFC3339))
	return rolloutIf.Patch(name, types.MergePatchType, []byte(patch))
}
