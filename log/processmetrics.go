package log

import "time"

func ProcessMetrics(serviceName string, success bool, responseTime time.Duration) {
	getMetricsLog("ProcessMetrics", success, responseTime).
		Str("serviceName", serviceName).
		Msg("")
}
