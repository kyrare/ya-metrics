package client

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"net/http"
)

type Client struct {
	serverAddr string
}

func (c *Client) Send(metricType metrics.MetricType, metric string, value float64) error {
	response, err := http.Post(
		fmt.Sprintf("http://%s/update/%s/%s/%v", c.serverAddr, metricType, metric, value),
		"text/plain",
		nil,
	)

	if err != nil {
		return err
	}

	return response.Body.Close()
}

func NewClient(serverAddr string) *Client {
	return &Client{serverAddr: serverAddr}
}
