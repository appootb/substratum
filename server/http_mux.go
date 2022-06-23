package server

import (
	"net/http"
	"runtime/debug"

	pctx "github.com/appootb/substratum/v2/plugin/context"
)

type httpServeMux struct {
	component string
	serveMux  *http.ServeMux
}

func (h *httpServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if result := recover(); result != nil {
			debug.PrintStack()
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	ctx := pctx.WithImplementContext(r.Context(), h.component)
	h.serveMux.ServeHTTP(w, r.WithContext(ctx))
}
