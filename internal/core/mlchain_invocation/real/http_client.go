package real

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/mlchain/mlchain-plugin-daemon/internal/core/mlchain_invocation"
)

func NewMlchainInvocationDaemon(base string, calling_key string) (mlchain_invocation.BackwardsInvocation, error) {
	var err error
	invocation := &RealBackwardsInvocation{}
	baseurl, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 120 * time.Second,
			}).Dial,
			IdleConnTimeout: 120 * time.Second,
		},
	}

	invocation.mlchainInnerApiBaseurl = baseurl
	invocation.client = client
	invocation.mlchainInnerApiKey = calling_key

	return invocation, nil
}
