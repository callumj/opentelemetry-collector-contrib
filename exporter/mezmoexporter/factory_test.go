// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package mezmoexporter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/exportertest"
)

func TestType(t *testing.T) {
	factory := NewFactory()
	pType := factory.Type()
	assert.Equal(t, pType, component.Type(typeStr))
}

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	assert.Equal(t, cfg, &Config{
		IngestURL: defaultIngestURL,
		IngestKey: "",

		HTTPClientSettings: confighttp.HTTPClientSettings{
			Timeout: 5 * time.Second,
		},
		RetrySettings: exporterhelper.NewDefaultRetrySettings(),
		QueueSettings: exporterhelper.NewDefaultQueueSettings(),
	})
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestIngestUrlMustConform(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.IngestURL = "/collector"
	cfg.IngestKey = "1234-1234"

	assert.Error(t, cfg.Validate(), `"ingest_url" must contain a valid host`)
}

func TestCreateLogsExporter(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.IngestURL = "https://example.com:8088/otel/ingest/rest"
	cfg.IngestKey = "1234-1234"

	params := exportertest.NewNopCreateSettings()
	_, err := createLogsExporter(context.Background(), params, cfg)
	assert.NoError(t, err)
}

func TestCreateLogsExporterNoConfig(t *testing.T) {
	params := exportertest.NewNopCreateSettings()
	_, err := createLogsExporter(context.Background(), params, nil)
	assert.Error(t, err)
}
