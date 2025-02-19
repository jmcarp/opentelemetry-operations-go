// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type selectorCloudMonitoring struct{}

var _ export.AggregatorSelector = selectorCloudMonitoring{}

// NewWithCloudMonitoringDistribution return a simple aggregation selector
// that uses lastvalue, counter, array, and aggregator for three kinds of metric.
//
// NOTE: this selector is just to ensure that LastValue is used for
// GaugeKind and HistogramKind.
// All other metric kinds have Sum default aggregation
// c.f. https://github.com/open-telemetry/opentelemetry-go/blob/main/sdk/metric/selector/simple/simple.go
//
// TODO: Remove this once SDK implements such a
// selector, otherwise Views API gives flexibility to set aggregation type on
// configuring measurement.
// c.f. https://github.com/open-telemetry/oteps/pull/89
func NewWithCloudMonitoringDistribution() export.AggregatorSelector {
	return selectorCloudMonitoring{}
}

func (selectorCloudMonitoring) AggregatorFor(descriptor *sdkapi.Descriptor, aggPtrs ...*aggregator.Aggregator) {
	switch descriptor.InstrumentKind() {
	case sdkapi.GaugeObserverInstrumentKind, sdkapi.HistogramInstrumentKind:
		aggs := lastvalue.New(len(aggPtrs))
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		aggs := sum.New(len(aggPtrs))
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	}
}
