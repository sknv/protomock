package render

import (
	"io"

	"github.com/goccy/go-json"
)

func DecodeJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v) //nolint:wrapcheck
}
