package provider

import "time"

// IFetcher can fetch data in streaming
type IFetcher interface {
	Fetch() <-chan []byte
	Size() int
	Start()
	Stop()
}

type IFetcherV2 interface {
	Fetch() <-chan interface{}
	Size() int
	Start()
	Stop()
}

type IFetchHandler interface {
	Name() string
	Size() int
	Action() (interface{}, error)
	Interval() time.Duration
	Close() error
}
