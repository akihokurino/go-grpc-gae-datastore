package government

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"gae-go-sample/adapter"

	"github.com/pkg/errors"
)

func NewClient(key string, baseURL *url.URL) adapter.GovernmentClient {
	return &client{
		key:     key,
		baseURL: baseURL,
	}
}

type client struct {
	key     string
	baseURL *url.URL
}

func (c *client) Get(ctx context.Context, pathStr string, params map[string]string) ([]byte, error) {
	buildURL, _ := url.Parse(c.baseURL.String())
	buildURL.Path = path.Join(buildURL.Path, pathStr)

	query := buildURL.Query()

	for k, v := range params {
		query.Set(k, v)
	}

	buildURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", buildURL.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("X-API-KEY", c.key)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return byteArray, nil
}
