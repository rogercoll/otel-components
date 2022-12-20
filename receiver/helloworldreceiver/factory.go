package helloreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

const (
	typeStr   = "hello_stats"
	stability = component.StabilityLevelDevelopment
)

func NewFactory() component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		createDefaultReceiverConfig,
		component.WithMetricsReceiver(createMetricsReceiver, stability))
}

func createDefaultConfig() *scraperhelper.ScraperControllerSettings {
	return &scraperhelper.ScraperControllerSettings{
		ReceiverSettings:   config.NewReceiverSettings(component.NewID(typeStr)),
		CollectionInterval: 10 * time.Second,
	}
}

func createDefaultReceiverConfig() component.Config {
	return createDefaultConfig()
}

func createMetricsReceiver(
	ctx context.Context,
	params component.ReceiverCreateSettings,
	config component.Config,
	consumer consumer.Metrics,
) (component.MetricsReceiver, error) {
	scConf := config.(*scraperhelper.ScraperControllerSettings)
	dsr, err := newReceiver(*scConf, params, consumer)
	if err != nil {
		return nil, err
	}

	return dsr, nil
}
