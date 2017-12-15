package fetcher

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// RSFetcher is a small interface definning the main act of fetching resources.
// This can be overwritten by the user of the client to provide more custom fetch behavior
type OAIFetcher interface {
	Fetch(source string) ([]byte, int, error)
}

// ErrNon200Response is returned from the BasicRSFetcher for any non-200 response
var ErrNon200Response = errors.New("Non-200 status code returned")

// Config is a simple implementation of the Fetcher interface. It is safe to use
// but limited in capability. No timeouts or extra headers are defined. The general recommendation
// is that the user of this client write their own implementation of the Fetcher interface.
type Config struct {
	Timeout time.Duration
}

type basicOAIFetcher struct {
	client  *http.Client
	timeout time.Duration
}

// NewBasicOAIFetcher instatiates a new basic fetcher for use in retrieving OAI links
func NewBasicOAIFetcher(c Config) *basicOAIFetcher {
	if c.Timeout == 0 {
		c.Timeout = time.Second * 600 // 5 minute timeout
	}

	return &basicOAIFetcher{
		client:  &http.Client{},
		timeout: c.Timeout,
	}
}

// Fetch retrieves the resource from source and writes it to dest. It is the callers responsibility
// to clear up any local files when they are finished with.
// This fetcher implementation will return an error for a non-200 response.
func (f *basicOAIFetcher) Fetch(source string) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel() // frees up the context resources.

	req, _ := http.NewRequest("GET", source, nil)
	// req.Close is needed to force not re-using connections
	//(see https://code.google.com/p/go/issues/detail?id=4677)
	req.Close = true
	req = req.WithContext(ctx)

	response, err := f.client.Do(req)
	if err != nil {
		// If we encountered a context.DeadlineExpired as part of making the initial connection, unlikely but possible,
		// this condition ensures we return the correct error type and track the correct stat.
		if ctx.Err() == context.DeadlineExceeded {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}
	defer response.Body.Close()

	statusGroup := response.StatusCode / 100
	switch statusGroup {
	case 2:
		return f.handle2xx(ctx, response, source)
	case 3:
		// The http.Client will follow redirects of its own accord. This handles other 3xx cases
		return nil, response.StatusCode, ErrNon200Response
	case 4:
		return nil, response.StatusCode, ErrNon200Response
	case 5:
		return nil, response.StatusCode, ErrNon200Response
	default:
		return nil, response.StatusCode, ErrNon200Response
	}
}

func (f *basicOAIFetcher) handle2xx(ctx context.Context, response *http.Response, source string) ([]byte, int, error) {
	if response.StatusCode != http.StatusOK {
		// return an err as we should only be getting 200
		return nil, response.StatusCode, ErrNon200Response
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}
	return data, response.StatusCode, nil
}
