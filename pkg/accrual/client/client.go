package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	Client struct {
		baseURL string
		http    *http.Client
		log     *zap.SugaredLogger
	}

	RequestOptions struct {
		Name   string
		Path   string
		Result interface{}
	}

	RequestError struct {
		Err          error
		StatusCode   int
		ResponseBody []byte
	}
)

func (e *RequestError) Error() string {
	return e.Err.Error()
}

func New(config *config.Config) *Client {
	return &Client{
		baseURL: config.AccrualSystemAddress,
		http:    &http.Client{},
	}
}

func (c *Client) get(r RequestOptions) (string, error) {
	return c.doRequest(http.MethodGet, r)
}

func (c *Client) doRequest(method string, r RequestOptions) (string, error) {
	req, err := http.NewRequest(method, c.baseURL+r.Path, nil)
	if err != nil {
		return "", &RequestError{
			Err: errors.Wrap(err, "build request"),
		}
	}

	var res *http.Response
	start := time.Now()
	defer func() {
		c.logRequest(r.Name, req.Method, start, res)
	}()

	res, err = c.http.Do(req)
	if err != nil {
		return "", &RequestError{
			Err: errors.Wrapf(err, "send request [%s]", req.URL),
		}
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", &RequestError{
			Err: errors.Wrapf(err, "read response body [%s]", req.URL),
		}
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", &RequestError{
			Err:          fmt.Errorf("bad response code for [%s] %d: %s", req.URL, res.StatusCode, resBody),
			StatusCode:   res.StatusCode,
			ResponseBody: resBody,
		}
	}

	if err := json.Unmarshal(resBody, r.Result); err != nil {
		return "", &RequestError{
			Err:          errors.Wrapf(err, "unmarshal response body [%s]", req.URL),
			StatusCode:   res.StatusCode,
			ResponseBody: resBody,
		}
	}

	return string(resBody), nil
}

func (c *Client) logRequest(name, method string, start time.Time, res *http.Response) {
	var status int
	if res != nil {
		status = res.StatusCode
	}

	c.log.With(
		"name", name,
		"method", method,
		"duration", time.Since(start).Seconds(),
		"status_code", status,
	).Info("accrual system request")
}
