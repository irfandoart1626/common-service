package log

import "time"

func BackendMetrics(backend, serviceName string, success bool, responseTime time.Duration) {
	getMetricsLog("BackendMetrics", success, responseTime).
		Str("backend", backend).
		Str("serviceName", serviceName).
		Msg("")
}
