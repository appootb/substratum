package golang

const serviceTpl = `
var {{ subjectsName . }} = map[string][]permission.Subject{
	{{- range $url, $subs := (serviceSubjects .) }}
	"{{ $url }}": {
		{{ range $subs }}permission.Subject_{{ . }},{{- end }}
	},
	{{- end }}
}

var {{ rolesName . }} = map[string][]string{
	{{- range $url, $roles := (serviceRoles .) }}
	"{{ $url }}": {
		{{ range $roles }}"{{ . }}",{{- end }}
	},
	{{- end }}
}

{{- if hasWebApi . }}
	{{ template "wrapper" . }}
{{- end }}

// Register scoped server.
func Register{{ serviceName . }}ScopeServer(component string, auth service.Authenticator, impl service.Implementor, srv {{ serverName . }}) error {
	// Register service required subjects.
	auth.RegisterServiceSubjects(component, {{ subjectsName . }}, {{ rolesName . }})

	// Register scoped gRPC server.
	for _, gRPC := range impl.GetGRPCServer(permission.VisibleScope_{{ (serviceScope .) }}) {
		Register{{ serverName . }}(gRPC, srv)
	}

	{{- if hasWebApi . }}
	// Register scoped gateway handler server.
	wrapper := wrapper{{ serverName . }}{
		{{ serverName . }}: srv,
		Implementor: impl,
	}
	{{- end }}

	{{- if or (hasWebApi .) (hasWebStream .) }}
	for _, mux := range impl.GetGatewayMux(permission.VisibleScope_{{ (serviceScope .) }}) {
		{{- if hasWebApi . }}
		// Register gateway handler.
		if err := Register{{ serviceName . }}HandlerServer(impl.Context(), mux, &wrapper); err != nil {
			return err
		}
		{{- end }}
		{{- if hasWebStream . }}
		// Register websocket handler.
		if err := Register{{ serviceName . }}WsHandlerServer(mux, srv, impl.StreamInterceptor()); err != nil {
			return err
		}
		{{- end }}
	}
	{{- end }}
	return nil
}

{{ if hasWebStream . }}
	{{ template "streaming" . }}
{{ end }}
`
