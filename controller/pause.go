package controller

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	logutil "github.com/argoproj/argo-rollouts/utils/log"
)

func completedPauseStep(rollout *v1alpha1.Rollout, pause *v1alpha1.RolloutPause) bool {
	logCtx := logutil.WithRollout(rollout)

	if pause != nil && pause.Duration != nil {
		now := metav1.Now()
		if rollout.Status.PauseStartTime != nil {
			expiredTime := rollout.Status.PauseStartTime.Add(time.Duration(*pause.Duration) * time.Second)
			if now.After(expiredTime) {
				logCtx.Info("Rollout has waited the duration of the pause step")
				return true
			}
		}
	}
	if pause != nil && pause.Duration == nil && rollout.Status.PauseStartTime != nil && !rollout.Spec.Paused {
		logCtx.Info("Rollout has been unpaused")
		return true
	}
	return false
}

func (c *Controller) checkEnqueueRolloutDuringPause(rollout *v1alpha1.Rollout, pause v1alpha1.RolloutPause) {
	logCtx := logutil.WithRollout(rollout)
	if pause.Duration == nil {
		return
	}
	if rollout.Status.PauseStartTime == nil {
		return
	}
	now := metav1.Now()
	expiredTime := rollout.Status.PauseStartTime.Add(time.Duration(*pause.Duration) * time.Second)
	nextResync := now.Add(c.resyncPeriod)
	if nextResync.After(expiredTime) && expiredTime.After(now.Time) {
		timeRemaining := expiredTime.Sub(now.Time)
		logCtx.Infof("Enqueueing Rollout in %s seconds", timeRemaining.String())
		c.enqueueRolloutAfter(rollout, timeRemaining)
	}
}

// calculatePauseStatus finds the fields related to a pause step for a rollout. If the pause is nil,
// the rollout will use the previous values
func calculatePauseStatus(rollout *v1alpha1.Rollout, addPause bool) (*metav1.Time, bool) {
	logCtx := logutil.WithRollout(rollout)
	pauseStartTime := rollout.Status.PauseStartTime
	paused := rollout.Spec.Paused
	if !paused {
		pauseStartTime = nil
	}
	if addPause {
		if pauseStartTime == nil {
			now := metav1.Now()
			logCtx.Infof("Setting PauseStartTime to %s", now.UTC().Format(time.RFC3339))
			pauseStartTime = &now
			paused = true
		}
	}
	return pauseStartTime, paused
}
