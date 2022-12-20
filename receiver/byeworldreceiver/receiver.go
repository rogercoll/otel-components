package byereceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

func newReceiver(config scraperhelper.ScraperControllerSettings, set component.ReceiverCreateSettings, nextConsumer consumer.Metrics) (component.MetricsReceiver, error) {
	scrp, err := scraperhelper.NewScraper("byeworld", scrape)
	if err != nil {
		return nil, err
	}
	return scraperhelper.NewScraperControllerReceiver(&config, set, nextConsumer, scraperhelper.AddScraper(scrp))
}

func greeterMetrics(name string) pmetric.Metrics {
	md := pmetric.NewMetrics()

	rs := md.ResourceMetrics().AppendEmpty()

	resourceAttr := rs.Resource().Attributes()
	resourceAttr.PutStr("greeter.name", name)

	ms := rs.ScopeMetrics().AppendEmpty().Metrics()
	m := ms.AppendEmpty()
	m.SetName("bye.requests")
	m.SetUnit("requests")
	m.SetEmptyGauge().DataPoints().AppendEmpty().SetIntValue(1)

	return md
}

func scrape(ctx context.Context) (pmetric.Metrics, error) {
	md := pmetric.NewMetrics()

	greeterMetrics("bob").ResourceMetrics().CopyTo(md.ResourceMetrics())

	return md, nil
}
