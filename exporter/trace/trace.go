// Copyright 2019 OpenTelemetry Authors
// Copyright 2021 Google LLC
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

package trace

import (
	"context"
	"fmt"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	traceclient "cloud.google.com/go/trace/apiv2"
	"google.golang.org/api/option"
	tracepb "google.golang.org/genproto/googleapis/devtools/cloudtrace/v2"
)

// traceExporter is an implementation of trace.Exporter and trace.BatchExporter
// that uploads spans to Stackdriver Trace in batch.
type traceExporter struct {
	o *options
	// uploadFn defaults in uploadSpans; it can be replaced for tests.
	uploadFn  func(ctx context.Context, spans []*tracepb.Span) error
	client    *traceclient.Client
	projectID string
	overflowLogger
}

func newTraceExporter(o *options) (*traceExporter, error) {
	clientOps := append(o.traceClientOptions, option.WithUserAgent(userAgent))
	client, err := traceclient.NewClient(o.context, clientOps...)
	if err != nil {
		return nil, fmt.Errorf("stackdriver: couldn't initiate trace client: %v", err)
	}
	e := &traceExporter{
		projectID:      o.projectID,
		client:         client,
		o:              o,
		overflowLogger: overflowLogger{delayDur: 5 * time.Second},
	}
	e.uploadFn = e.uploadSpans
	return e, nil
}

func (e *traceExporter) ExportSpans(ctx context.Context, spanData []sdktrace.ReadOnlySpan) error {
	// Ship the whole bundle o data.
	results := make([]*tracepb.Span, len(spanData))
	for i, sd := range spanData {
		results[i] = e.ConvertSpan(ctx, sd)
	}
	return e.uploadFn(ctx, results)
}

// ConvertSpan converts a ReadOnlySpan to Stackdriver Trace.
func (e *traceExporter) ConvertSpan(_ context.Context, sd sdktrace.ReadOnlySpan) *tracepb.Span {
	return e.protoFromReadOnlySpan(sd, e.projectID)
}

func (e *traceExporter) Shutdown(ctx context.Context) error {
	return e.client.Close()
}

// uploadSpans sends a set of spans to Stackdriver.
func (e *traceExporter) uploadSpans(ctx context.Context, spans []*tracepb.Span) error {
	req := tracepb.BatchWriteSpansRequest{
		Name:  "projects/" + e.projectID,
		Spans: spans,
	}

	var cancel func()
	ctx, cancel = newContextWithTimeout(ctx, e.o.timeout)
	defer cancel()

	// TODO(ymotongpoo): add this part after OTel support NeverSampler
	// for tracer.Start() initialization.
	//
	// tracer := apitrace.Register()
	// ctx, span := tracer.Start(
	// 	ctx,
	// 	"go.opentelemetry.io/otel/exporters/stackdriver.uploadSpans",
	// )
	// defer span.End()
	// span.SetAttributes(kv.Int64("num_spans", int64(len(spans))))

	err := e.client.BatchWriteSpans(ctx, &req)
	if err != nil {
		// TODO(ymotongpoo): handle detailed error categories
		// span.SetStatus(codes.Unknown)
		e.o.handleError(fmt.Errorf("failed to export to Google Cloud Trace: %w", err))
	}
	return err
}

// overflowLogger ensures that at most one overflow error log message is
// written every 5 seconds.
type overflowLogger struct {
	delayDur time.Duration
}
