// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package filestatsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filestatsreceiver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filestatsreceiver/internal/metadata"
)

// NewFactory creates a new filestats receiver factory.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		newDefaultConfig,
		receiver.WithMetrics(newReceiver, metadata.MetricsStability))
}

func newDefaultConfig() component.Config {
	return &Config{
		ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
			CollectionInterval: 1 * time.Minute,
		},
		MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
	}
}

func newReceiver(
	_ context.Context,
	settings receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	fileStatsConfig := cfg.(*Config)
	metricsBuilder := metadata.NewMetricsBuilder(fileStatsConfig.MetricsBuilderConfig, settings)

	mp := newScraper(metricsBuilder, fileStatsConfig, settings.TelemetrySettings.Logger)
	s, err := scraperhelper.NewScraper(settings.ID.Name(), mp.scrape)
	if err != nil {
		return nil, err
	}
	opt := scraperhelper.AddScraper(s)

	return scraperhelper.NewScraperControllerReceiver(
		&fileStatsConfig.ScraperControllerSettings,
		settings,
		consumer,
		opt,
	)
}
