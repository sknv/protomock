package middleware

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func HandleError(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		err := next(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return err // Return the err in case there other middlewares.
	}
}
