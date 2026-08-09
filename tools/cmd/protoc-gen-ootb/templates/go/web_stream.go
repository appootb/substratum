package golang

const webStreamTpl = `
func Register{{ serviceName . }}WsHandlerServer(mux *runtime.ServeMux, srv {{ serverName . }}, streamInterceptor grpc.StreamServerInterceptor) error {
	{{ range .Methods }}
		{{- if isWebStream . }}
			{{- if .ClientStreaming }}
				{{ template "clientStream" . }}
			{{- else if .ServerStreaming }}
				{{ template "serverStream" . }}
			{{- end }}
		{{- end }}
	{{ end }}
	return nil
}
`
