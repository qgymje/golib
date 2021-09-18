package httputil

import (
	"fmt"
	"net/http"
	"time"
)

func CreateHttpClient(timeout time.Duration, idle, perHost int) *http.Client {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer
	defaultTransport.MaxIdleConns = idle
	defaultTransport.MaxIdleConnsPerHost = perHost

	client := &http.Client{}

	client.Transport = &defaultTransport
	client.Timeout = timeout

	return client
}
