package http

import (
	"fmt"
	"net/http"

	"github.com/sknv/protomock/pkg/http/render"
)

type (
	ReponseHeaders map[string]string
	ReponseBody    map[string]any
)

type Response struct {
	Status  int            `json:"status"`
	Headers ReponseHeaders `json:"headers"`
	Body    ReponseBody    `json:"body"`
}

func (r Response) JSON(w http.ResponseWriter) error {
	for header, value := range r.Headers {
		w.Header().Set(header, value)
	}

	if err := render.JSON(w, r.Status, r.Body); err != nil {
		return fmt.Errorf("render json: %w", err)
	}

	return nil
}
