package log

import (
	"time"

	"github.com/rs/zerolog"
)

func getMetricsLog(metricsType string, success bool, responseTime time.Duration) *zerolog.Event {
	successRate := 0
	errorCount := 1
	if success {
		successRate = 100
		errorCount = 0
	}
	messageCount := 1
	l := GetLogger()
	return l.Log().Str("type", metricsType).
		Int("messageCount", messageCount).
		Int("errorCount", errorCount).
		Int("successRate", successRate).
		Int64("responseTimeInMs", responseTime.Milliseconds())

}
