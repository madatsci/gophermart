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
		RetryAfter   time.Duration
	}
)

func (e *RequestError) Error() string {
	return e.Err.Error()
}

func New(config *config.Config, logger *zap.SugaredLogger) *Client {
	return &Client{
		baseURL: config.AccrualSystemAddress,
		http:    &http.Client{},
		log:     logger,
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

	if res.StatusCode != http.StatusOK {
		reqErr := &RequestError{
			Err:          fmt.Errorf("bad response code for [%s] %d: %s", req.URL, res.StatusCode, resBody),
			StatusCode:   res.StatusCode,
			ResponseBody: resBody,
		}
		if res.StatusCode == http.StatusTooManyRequests {
			if durationStr := res.Header.Get("Retry-After"); durationStr != "" {
				duration, err := time.ParseDuration(durationStr + "s")
				if err != nil {
					c.log.Errorf("could not parse duration from response header: %s", durationStr)
				} else {
					reqErr.RetryAfter = duration
				}
			}
		}

		return "", reqErr
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
