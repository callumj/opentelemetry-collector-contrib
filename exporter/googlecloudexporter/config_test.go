// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package googlecloudexporter

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector"
	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	sub, err := cm.Sub(component.NewIDWithName(typeStr, "").String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))

	assert.Equal(t, sanitize(cfg.(*Config)), sanitize(factory.CreateDefaultConfig().(*Config)))

	sub, err = cm.Sub(component.NewIDWithName(typeStr, "customname").String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))

	assert.Equal(t,
		&Config{
			TimeoutSettings: exporterhelper.TimeoutSettings{
				Timeout: 20 * time.Second,
			},
			Config: collector.Config{
				ProjectID: "my-project",
				UserAgent: "opentelemetry-collector-contrib {{version}}",
				LogConfig: collector.LogConfig{
					ServiceResourceLabels: true,
				},
				MetricConfig: collector.MetricConfig{
					Prefix:                           "prefix",
					SkipCreateMetricDescriptor:       true,
					KnownDomains:                     []string{"googleapis.com", "kubernetes.io", "istio.io", "knative.dev"},
					CreateMetricDescriptorBufferSize: 10,
					InstrumentationLibraryLabels:     true,
					ServiceResourceLabels:            true,
					ClientConfig: collector.ClientConfig{
						Endpoint:    "test-metric-endpoint",
						UseInsecure: true,
					},
				},
				TraceConfig: collector.TraceConfig{
					ClientConfig: collector.ClientConfig{
						Endpoint:    "test-trace-endpoint",
						UseInsecure: true,
					},
				},
			},
			RetrySettings: exporterhelper.RetrySettings{
				Enabled:             true,
				InitialInterval:     10 * time.Second,
				MaxInterval:         1 * time.Minute,
				MaxElapsedTime:      10 * time.Minute,
				RandomizationFactor: backoff.DefaultRandomizationFactor,
				Multiplier:          backoff.DefaultMultiplier,
			},
			QueueSettings: exporterhelper.QueueSettings{
				Enabled:      true,
				NumConsumers: 2,
				QueueSize:    10,
			},
		},
		sanitize(cfg.(*Config)))
}

func sanitize(cfg *Config) *Config {
	cfg.Config.MetricConfig.MapMonitoredResource = nil
	cfg.Config.MetricConfig.GetMetricName = nil
	return cfg
}
