package http

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/uptrace/bunrouter"

	"github.com/sknv/protomock/pkg/http/render"
)

type (
	MockRequestParams  map[string]string
	MockRequestHeaders map[string]string
	MockRequestBody    map[string]any
)

type MockRequest struct {
	Params  MockRequestParams  `json:"params"`
	Headers MockRequestHeaders `json:"headers"`
	Body    MockRequestBody    `json:"body"`
}

func NewMockRequestFrom(r bunrouter.Request) (MockRequest, error) {
	var body MockRequestBody
	if err := render.DecodeJSON(r.Body, &body); err != nil && !errors.Is(err, io.EOF) {
		return MockRequest{}, fmt.Errorf("decode json body: %w", err)
	}

	headers := make(MockRequestHeaders, len(r.Header))
	for header := range r.Header {
		headers[header] = strings.ToLower(r.Header.Get(header))
	}

	return MockRequest{
		Params:  r.Params().Map(),
		Headers: headers,
		Body:    body,
	}, nil
}
