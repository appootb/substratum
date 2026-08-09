package golang

const serverStreamTpl = `
	if err := mux.HandlePath("GET", "{{ webStreamURL . }}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		fn := func(c *websocket.Conn) {
			ctx, err := runtime.AnnotateIncomingContext(r.Context(), mux, r, "{{ fullMethod . }}")
			if err != nil {
				_ = c.WriteClose(http.StatusBadRequest)
				return
			}
			inbound, outbound := runtime.MarshalerForRequest(mux, r)
			stream := gateway.NewWebsocketStream(ctx, c, inbound, outbound)
			req := new({{ inputMessage . }})
			if err = stream.RecvMsg(req); err != nil {
				_ = c.WriteClose(http.StatusBadRequest)
				return
			}
			// protoc-gen-go-grpc v1.3+ uses GenericServerStream (no unexported *XServer wrappers).
			ss := &grpc.GenericServerStream[{{ inputMessage . }}, {{ outputMessage . }}]{ServerStream: stream}
			if streamInterceptor == nil {
				err = srv.{{ name . }}(req, ss)
				if err != nil {
					_ = c.WriteClose(http.StatusInternalServerError)
				}
				return
			}
			handler := func(_ interface{}, stream grpc.ServerStream) error {
				return srv.{{ name . }}(req, &grpc.GenericServerStream[{{ inputMessage . }}, {{ outputMessage . }}]{ServerStream: stream})
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
