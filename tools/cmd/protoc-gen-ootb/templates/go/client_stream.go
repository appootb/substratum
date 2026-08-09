package golang

const clientStreamTpl = `
	if err := mux.HandlePath("GET", "{{ webStreamURL . }}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		fn := func(c *websocket.Conn) {
			ctx, err := runtime.AnnotateIncomingContext(r.Context(), mux, r, "{{ fullMethod . }}")
			if err != nil {
				_ = c.WriteClose(http.StatusBadRequest)
				return
			}
			inbound, outbound := runtime.MarshalerForRequest(mux, r)
			stream := gateway.NewWebsocketStream(ctx, c, inbound, outbound)
			if streamInterceptor == nil {
				err = srv.{{ name . }}(&{{ streamServerImpl . }}{stream})
				if err != nil {
					_ = c.WriteClose(http.StatusInternalServerError)
				}
				return
			}
			handler := func(_ interface{}, stream grpc.ServerStream) error {
				return srv.{{ name . }}(&{{ streamServerImpl . }}{stream})
			}
			info := &grpc.StreamServerInfo{
				FullMethod:     "{{ fullMethod . }}",
				IsClientStream: {{ .ClientStreaming }},
				IsServerStream: {{ .ServerStreaming }},
			}
			if err = streamInterceptor(srv, stream, info, handler); err != nil {
				_ = c.WriteClose(http.StatusInternalServerError)
			}
		}
		websocket.Handler(fn).ServeHTTP(w, r)
	}); err != nil {
		return err
	}
`
