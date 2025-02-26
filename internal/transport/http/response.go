package http

import (
	"fmt"
	"net/http"

	"github.com/sknv/protomock/pkg/http/render"
)

type (
	MockResponseBody map[string]any
)

type MockResponse struct {
	Status int              `json:"status"`
	Body   MockResponseBody `json:"body"`
}

func (r MockResponse) JSON(w http.ResponseWriter) error {
	if err := render.JSON(w, r.Status, r.Body); err != nil {
		return fmt.Errorf("render json: %w", err)
	}

	return nil
}
