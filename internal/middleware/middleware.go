package middleware

import (
	"net/http"

	"github.com/bbarrington0099/Chirpy/internal/apiConfig"
)

func (api *apiConfig.Conf) middlewareFileserverHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}