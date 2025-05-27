package middleware

import (
	"net/http"

	"github.com/bbarrington0099/Chirpy/internal/apiconfig"
)

type LocalConf apiconfig.Conf

func (api *LocalConf) MiddlewareFileserverHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(*apiconfig.Conf)(api).FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}