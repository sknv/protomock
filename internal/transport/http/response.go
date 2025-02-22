package http

import (
	"fmt"
	"net/http"

	"github.com/sknv/protomock/pkg/http/render"
)

type (
	MockResponseHeaders map[string]string
	MockResponseBody    map[string]any
)

type MockResponse struct {
	Status  int                 `json:"status"`
	Headers MockResponseHeaders `json:"headers"`
	Body    MockResponseBody    `json:"body"`
}

func (r MockResponse) JSON(w http.ResponseWriter) error {
	for header, value := range r.Headers {
		w.Header().Set(header, value)
	}

	if err := render.JSON(w, r.Status, r.Body); err != nil {
		return fmt.Errorf("render json: %w", err)
	}

	return nil
}
