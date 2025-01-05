package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	Interval int64     `json:"interval"`
	Requests []Request `json:"requests"`
}

func NewClient(interval int64, requests []Request) *Client {
	return &Client{
		Interval: interval,
		Requests: requests,
	}

}

type Request struct {
	URI     string `json:"uri"`
	Method  string `json:"method"`
	Payload string `json:"payload"`
}

func (c *Client) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(c.Interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			var wg sync.WaitGroup
			for _, req := range c.Requests {
				wg.Add(1)
				go func() {
					defer wg.Done()
					req.Do()
				}()

			}
			wg.Wait()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (r *Request) Do() {
	req, err := http.NewRequest(r.Method, r.URI, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return
	}
	fmt.Printf("client: response body: %s\n", resBody)
}
