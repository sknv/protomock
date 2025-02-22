package http

import (
	"fmt"
	"net/http"

	"github.com/uptrace/bunrouter"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Route(router *bunrouter.Router) {
	router.GET("/users", h.handleHttpRequest)
	router.GET("/users/:id", h.handleHttpRequest)
}

func (h *Handlers) handleHttpRequest(w http.ResponseWriter, r bunrouter.Request) error {
	fmt.Println(r.Method, r.Route(), r.Params().Map())

	return nil
}
