// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	v1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	mock "github.com/stretchr/testify/mock"
)

// Provider is an autogenerated mock type for the Provider type
type Provider struct {
	mock.Mock
}

// GarbageCollect provides a mock function with given fields: _a0, _a1, _a2
func (_m *Provider) GarbageCollect(_a0 *v1alpha1.AnalysisRun, _a1 v1alpha1.Metric, _a2 int) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1alpha1.AnalysisRun, v1alpha1.Metric, int) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMetadata provides a mock function with given fields: metric
func (_m *Provider) GetMetadata(metric v1alpha1.Metric) map[string]string {
	ret := _m.Called(metric)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(v1alpha1.Metric) map[string]string); ok {
		r0 = rf(metric)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// Resume provides a mock function with given fields: _a0, _a1, _a2
func (_m *Provider) Resume(_a0 *v1alpha1.AnalysisRun, _a1 v1alpha1.Metric, _a2 v1alpha1.Measurement) v1alpha1.Measurement {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 v1alpha1.Measurement
	if rf, ok := ret.Get(0).(func(*v1alpha1.AnalysisRun, v1alpha1.Metric, v1alpha1.Measurement) v1alpha1.Measurement); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(v1alpha1.Measurement)
	}

	return r0
}

// Run provides a mock function with given fields: _a0, _a1
func (_m *Provider) Run(_a0 *v1alpha1.AnalysisRun, _a1 v1alpha1.Metric) v1alpha1.Measurement {
	ret := _m.Called(_a0, _a1)

	var r0 v1alpha1.Measurement
	if rf, ok := ret.Get(0).(func(*v1alpha1.AnalysisRun, v1alpha1.Metric) v1alpha1.Measurement); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(v1alpha1.Measurement)
	}

	return r0
}

// Terminate provides a mock function with given fields: _a0, _a1, _a2
func (_m *Provider) Terminate(_a0 *v1alpha1.AnalysisRun, _a1 v1alpha1.Metric, _a2 v1alpha1.Measurement) v1alpha1.Measurement {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 v1alpha1.Measurement
	if rf, ok := ret.Get(0).(func(*v1alpha1.AnalysisRun, v1alpha1.Metric, v1alpha1.Measurement) v1alpha1.Measurement); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(v1alpha1.Measurement)
	}

	return r0
}

// Type provides a mock function with given fields:
func (_m *Provider) Type() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewProvider creates a new instance of Provider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProvider(t mockConstructorTestingTNewProvider) *Provider {
	mock := &Provider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
