package golang

const wrapperTpl = `
type wrapper{{ serverName . }} struct {
	{{ serverName . }}
	service.Implementor
}

{{ range .Methods }}
	{{ template "method" . }}
{{ end }}
`
