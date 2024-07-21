package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"github.com/kyrare/ya-metrics/internal/domain/utils"
	"github.com/kyrare/ya-metrics/internal/infrastructure/compress"
	"go.uber.org/zap"
)

type Client struct {
	serverAddr string
	addKey     bool
	jobs       chan []metrics.Metrics
	Logger     zap.SugaredLogger
}

func (c *Client) Send(data []metrics.Metrics) {
	c.jobs <- data
}

func (c *Client) newWorker(jobs <-chan []metrics.Metrics) {
	c.Logger.Infoln("Создан новый воркер")

	for data := range jobs {
		err := c.send(data)
		if err != nil {
			c.Logger.Error(err)
		}
	}
}

func (c *Client) send(data []metrics.Metrics) error {
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

	bodyJSONCompress, err := compress.Compress(bodyJSON)

	if err != nil {
		c.Logger.Error("Произошла ошибка сжатия данных", err)
		return err
	}

	body := bytes.NewBuffer(bodyJSONCompress)

	req, err := http.NewRequest("POST", uri, body)

	if err != nil {
		c.Logger.Error("Не удалось создать реквест")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	if c.addKey {
		req.Header.Set("HashSHA256", utils.Hash(bodyJSON))
	}

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

func NewClient(serverAddr string, addKey bool, workersCount uint64, logger zap.SugaredLogger) *Client {
	jobs := make(chan []metrics.Metrics, workersCount)

	c := &Client{
		serverAddr: serverAddr,
		addKey:     addKey,
		jobs:       jobs,
		Logger:     logger,
	}

	var i uint64
	for i = 0; i < workersCount; i++ {
		go c.newWorker(jobs)
	}

	return c
}
