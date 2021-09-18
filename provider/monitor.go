package provider

import "time"

// IMonitor provides monitor
type IMonitor interface {
	Interval() time.Duration
	Send(v1, v2, v3 string, num int64)
}
