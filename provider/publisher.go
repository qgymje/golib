package provider

import "time"

type IPublisher interface {
	Publish(topic string, data interface{}, delay time.Duration) error
}
