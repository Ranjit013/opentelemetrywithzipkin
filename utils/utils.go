package utils

import (
	"bytes"
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"time"
)

func SendRequest(ctx context.Context, method string, url string, data []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	client := http.Client{
		// Wrap the Transport with one that starts a span and injects the span context
		// into the outbound request headers.
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}

	return client.Do(request)
}
