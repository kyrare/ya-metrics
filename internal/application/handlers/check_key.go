package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/kyrare/ya-metrics/internal/domain/utils"
)

func (h *Handler) checkRequestKey(r *http.Request) bool {
	if !h.checkKey {
		return true
	}

	headerHash := r.Header.Get("HashSHA256")

	fmt.Println("Send hash - ", headerHash)

	if headerHash == "" {
		return false
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	err := r.Body.Close()
	if err != nil {
		h.logger.Error(err)
		return false
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	fmt.Println("Body hash - ", utils.Hash(bodyBytes))

	return headerHash == utils.Hash(bodyBytes)
}
