// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

type finishOption interface {
	FinishEnabled() bool
}

func newExperimentalMeterFactory(options []experimentalOption) (meterFactoryFn meterFactory, shutdown func()) {
	for _, option := range options {
		if finish, ok := option.(finishOption); ok && finish.FinishEnabled() {
			return newFinishMeterFactory()
		}
	}
	return defaultMeterFactory, nil
}
