package byereceiver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

func TestScrape(t *testing.T) {
	consumer := make(mockConsumer)

	config := scraperhelper.ScraperControllerSettings{
		CollectionInterval: 10 * time.Millisecond,
	}

	r, err := newReceiver(config, componenttest.NewNopReceiverCreateSettings(), consumer)
	require.NoError(t, err)
	assert.NotNil(t, r)

	assert.NoError(t, r.Start(context.Background(), componenttest.NewNopHost()))

	md := <-consumer
	assert.Equal(t, md.ResourceMetrics().Len(), 1)

	rsm := md.ResourceMetrics().At(0)

	attr, exists := rsm.Resource().Attributes().Get("greeter.name")
	assert.True(t, exists)
	assert.Equal(t, attr.Str(), "bob")

	assert.NoError(t, r.Shutdown(context.Background()))
}

type mockConsumer chan pmetric.Metrics

func (m mockConsumer) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{}
}

func (m mockConsumer) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	m <- md
	return nil
}