package dispatch

import (
	"net/http"

	"simulator/service/src/modules/engine/application/ports"
)

// httpDispatcher sends over HTTP using the standard client, since each request
// targets an arbitrary URL/method (the BaseURL-oriented gokit httpclient does not
// fit a per-request URL).
type httpDispatcher struct {
	client *http.Client
}

// mqttDispatcher sends over MQTT using the gokit mqtt client.
type mqttDispatcher struct{}

// registry maps a protocol id to its dispatcher.
type registry struct {
	byProtocol map[string]ports.Dispatcher
}
