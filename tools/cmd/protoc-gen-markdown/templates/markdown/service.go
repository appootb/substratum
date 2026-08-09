package markdown

const serviceTpl = `
<h2 id="{{ anchorName . }}">{{ headerTitle . }}</h2>

> {{ leadingComment .SourceCodeInfo }}

{{ range .Methods }}
{{ template "method" . }}
{{ end }}
`
