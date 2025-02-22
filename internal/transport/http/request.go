package http

import "github.com/uptrace/bunrouter"

type (
	MockRequestParams  map[string]string
	MockRequestHeaders map[string]string
)

type MockRequest struct {
	Params  MockRequestParams  `json:"params"`
	Headers MockRequestHeaders `json:"headers"`
}

func NewMockRequestFrom(r bunrouter.Request) MockRequest {
	headers := make(MockRequestHeaders, len(r.Header))
	for header := range r.Header {
		headers[header] = r.Header.Get(header)
	}

	return MockRequest{
		Params:  r.Params().Map(),
		Headers: headers,
	}
}
