package http

import (
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
		h.handleMockRequest(router, mock)
	}
}

func (h *Handlers) handleMockRequest(router *bunrouter.Router, mock Mock) {
	router.Handle(mock.Method, mock.Path, func(w http.ResponseWriter, r bunrouter.Request) error {
		return nil
	})
}
