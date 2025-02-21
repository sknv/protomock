package http

import (
	"net/http"

	"github.com/sknv/protomock/pkg/http/routegroup"
	"github.com/sknv/protomock/pkg/log"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Route(router *routegroup.Group) {
	router.HandleFunc("/", h.handleHttpRequest)
}

func (h *Handlers) handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.FromContext(ctx).InfoContext(ctx, "Here")
}
