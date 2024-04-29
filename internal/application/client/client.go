package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Client struct {
	serverAddr string
	Logger     zap.SugaredLogger
}

func (c *Client) Send(metricType metrics.MetricType, metric string, value float64) error {
	bodyData := metrics.NewMetrics(metricType, metric, value)

	bodyJson, err := json.Marshal(bodyData)

	if err != nil {
		c.Logger.Error("Не удалось сконвертировать body в json", bodyData)
		return err
	}

	uri := fmt.Sprintf("http://%s/update/", c.serverAddr)

	c.Logger.Infoln(
		"Начало выполнение запроса",
		"uri", uri,
		"body", string(bodyJson),
		"method", "POST",
	)

	body := bytes.NewBuffer(bodyJson)

	response, err := http.Post(uri, "application/json", body)

	if err != nil {
		c.Logger.Error("Произошла ошибка отправки данных", err, response, string(bodyJson))
		return err
	}

	c.Logger.Infoln(
		"Запрос выполнен",
		"uri", uri,
		"method", "POST",
		"status", response.Status,
		"size", bodySize(response.Body),
	)

	return response.Body.Close()
}

func bodySize(body io.ReadCloser) int {
	bytes, err := io.ReadAll(body)

	if err != nil {
		return 0
	}

	return len(bytes)
}

func NewClient(serverAddr string, logger zap.SugaredLogger) *Client {
	return &Client{serverAddr: serverAddr, Logger: logger}
}
