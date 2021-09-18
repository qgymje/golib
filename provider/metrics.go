package provider

// IMetrics is an expiremental interface for metrics
// the method now is only fulfill influxdb
type IMetrics interface {
	Measure(name string, tags map[string]string, fields map[string]interface{}) error
}
