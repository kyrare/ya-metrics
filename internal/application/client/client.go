package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"github.com/kyrare/ya-metrics/internal/infrastructure/compress"
	"go.uber.org/zap"
)

type Client struct {
	serverAddr string
	Logger     zap.SugaredLogger
}

func (c *Client) Send(data []metrics.Metrics) error {
	bodyJSON, err := json.Marshal(data)

	if err != nil {
		c.Logger.Error("Не удалось сконвертировать body в json", data)
		return err
	}

	uri := "http://" + c.serverAddr + "/updates/"

	c.Logger.Infoln(
		"Начало выполнение запроса",
		"uri", uri,
		"body", string(bodyJSON),
		"method", "POST",
	)

	bodyJSON, err = compress.Compress(bodyJSON)

	if err != nil {
		c.Logger.Error("Произошла ошибка сжатия данных", err)
		return err
	}

	body := bytes.NewBuffer(bodyJSON)

	req, err := http.NewRequest("POST", uri, body)

	if err != nil {
		c.Logger.Error("Не удалось создать реквест")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		c.Logger.Error("Произошла ошибка отправки данных", err, response, string(bodyJSON))
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
	defer body.Close()

	if err != nil {
		return 0
	}

	return len(bytes)
}

func NewClient(serverAddr string, logger zap.SugaredLogger) *Client {
	return &Client{serverAddr: serverAddr, Logger: logger}
}
