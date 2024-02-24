// go-musclemem provides an Go SDK for the musclemem-api
// based on github.com/google/go-github
package musclemem

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	Version = "v1.0.0"

	// TODO: update API version
	defaultAPIVersion = "2024-01-01"
	defaultBaseURL    = "https://api.musclemem.com/"
	defaultUserAgent  = "go-musclemem" + "/" + Version
	defaultMediaType  = "application/json"

	headerAPIVersion = "X-Musclemem-Api-Version"
	headerRetryAfter = "Retry-After"
)

type client struct {
	baseURL *url.URL
	apiKey  string

	Users     *userService
	Workouts  *workoutService
	Exercises *exerciseService
}

func newClient(baseURL string, apiKey string) (*client, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if apiKey == "" {
		return nil, errors.New("no api key provided")
	}

	c := &client{
		baseURL: url,
		apiKey:  apiKey,
	}

	c.Users = &userService{client: c}
	c.Workouts = &workoutService{client: c}
	c.Exercises = &exerciseService{client: c}

	return c, nil
}

type service struct {
	client *client
}

// TODO: extend with request options
type requestOption func(r *http.Request)

// send sends handles requests to the musclemem api
func (c *client) send(
	ctx context.Context,
	method string,
	path string,
	body interface{},
	opts ...requestOption,
) (*http.Response, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.baseURL)
	}

	url, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Accept", defaultMediaType)
	req.Header.Set("Content-Type", defaultMediaType)
	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set(headerAPIVersion, defaultAPIVersion)

	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	// TODO: use *http.Client instead
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
