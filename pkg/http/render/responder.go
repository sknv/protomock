package render

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

// JSON renders JSON data response with the provided status,
// automatically escaping HTML and setting the Content-Type as application/json.
func JSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
