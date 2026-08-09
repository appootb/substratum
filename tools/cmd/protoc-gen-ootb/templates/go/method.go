package golang

const methodTpl = `
{{ if .ClientStreaming }}
	func (w *wrapper{{ serverName .Service }}) {{ name . }}(srv {{ serverStream . }}) error {
		return w.{{ serverName .Service }}.{{ name . }}(srv)
	}
{{ else if .ServerStreaming }}
	func (w *wrapper{{ serverName .Service }}) {{ name . }}(req *{{ inputMessage . }}, srv {{ serverStream . }}) error {
		return w.{{ serverName .Service }}.{{ name . }}(req, srv)
	}
{{ else }}
	func (w *wrapper{{ serverName .Service }}) {{ name . }}(ctx context.Context, req *{{ inputMessage . }}) (*{{ outputMessage . }}, error) {
		if w.UnaryInterceptor() == nil {
			return w.{{ serverName .Service }}.{{ name . }}(ctx, req)
		}
		info := &grpc.UnaryServerInfo{
			Server:     w.{{ serverName .Service }},
			FullMethod: "{{ fullMethod . }}",
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return w.{{ serverName .Service }}.{{ name . }}(ctx, req.(*{{ inputMessage . }}))
		}
		resp, err := w.UnaryInterceptor()(ctx, req, info, handler)
		if err != nil {
			return nil, err
		}
		return resp.(*{{ outputMessage . }}), nil
	}
{{ end }}
`
