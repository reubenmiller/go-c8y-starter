package tasks

import (
	"bytes"
	"strings"

	"github.com/reubenmiller/go-c8y/pkg/c8y"

	"github.com/reubenmiller/go-c8y/pkg/microservice"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	//
	// measurements metrics per device
	//
	measurements = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "c8y_measurements_total",
			Help: "Total number of measurements received",
		},
		[]string{"source", "label"},
	)
)

type MeasurementMetrics struct{}

func NewMeasurementMetrics(ms *microservice.Microservice) *MeasurementMetrics {
	v := &MeasurementMetrics{}

	labels := strings.Split(ms.Config.GetString("measurement.labels"), ",")

	ms.SubscribeToNotifications(
		ms.WithServiceUserCredentials(),
		c8y.RealtimeMeasurements("*"),
		onMeasurementNotification(labels),
	)
	return v
}

func onMeasurementNotification(labels []string) func(*c8y.Message) {
	return func(msg *c8y.Message) {
		if msg.Payload.RealtimeAction == "CREATE" {
			for i := range labels {
				if bytes.Contains(msg.Payload.Data, []byte(labels[i])) {
					measurements.
						WithLabelValues(msg.Channel[14:], labels[i]).
						Inc()
					break
				}
			}
		}
	}
}
