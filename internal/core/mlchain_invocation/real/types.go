package real

import (
	"net/http"
	"net/url"
)

type RealBackwardsInvocation struct {
	mlchainInnerApiKey     string
	mlchainInnerApiBaseurl *url.URL
	client              *http.Client
}

type BaseBackwardsInvocationResponse[T any] struct {
	Data  *T     `json:"data,omitempty"`
	Error string `json:"error"`
}
