package server

import (
	"net/http"
	"runtime/debug"

	sctx "github.com/appootb/substratum/v2/context"
	"github.com/appootb/substratum/v2/service"
)

type handlerWrapper struct {
	component string
	handler   http.Handler
}

func (h *handlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if result := recover(); result != nil {
			debug.PrintStack()
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	//
	ctx := service.ContextWithServiceMethod(r.Context(), &service.Method{
		FullMethod:    r.URL.Path,
		IsHttpGateway: true,
	})
	h.handler.ServeHTTP(w, r.WithContext(sctx.WithServerContext(ctx, h.component)))
}

type httpServeMux struct {
	component string
	serveMux  *http.ServeMux
}

func (h *httpServeMux) Handle(pattern string, handler http.Handler) {
	h.serveMux.Handle(pattern, &handlerWrapper{
		component: h.component,
		handler:   handler,
	})
}

func (h *httpServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h.serveMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if result := recover(); result != nil {
				debug.PrintStack()
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		//
		ctx := service.ContextWithServiceMethod(r.Context(), &service.Method{
			FullMethod:    r.URL.Path,
			IsHttpGateway: true,
		})
		handler(w, r.WithContext(sctx.WithServerContext(ctx, h.component)))
	})
}
