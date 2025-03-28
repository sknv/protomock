package http

import (
	"fmt"
	"net/http"

	"github.com/uptrace/bunrouter"
)

type Handlers struct {
	mocks Mocks
}

func NewHandlers(mocks Mocks) *Handlers {
	return &Handlers{
		mocks: mocks,
	}
}

func (h *Handlers) Route(router *bunrouter.Router) {
	for _, mock := range h.mocks {
		handleMockRequest(router, mock)
	}
}

func handleMockRequest(router *bunrouter.Router, mock Mock) {
	router.Handle(mock.Method, mock.Path, func(w http.ResponseWriter, r bunrouter.Request) error {
		ctx := r.Context()

		request, err := NewMockRequestFrom(r)
		if err != nil {
			return fmt.Errorf("decode request: %w", err)
		}

		response, err := mock.Eval(ctx, request)
		if err != nil {
			return fmt.Errorf("evaluate mock: %w", err)
		}

		return response.JSON(w)
	})
}
