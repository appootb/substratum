package markdown

const messageTpl = `
{{ if .IsWellKnown }}

> JSON Demo

{{ jsonWellKnown . }}
{{ else }}
> {{ leadingComment .SourceCodeInfo }}

|Field|Proto Type|JSON Type|Validation|Comment|{{ if .Syntax.SupportsRequiredPrefix }}Default|Required|{{ end }}
|---|---|---|---|---|{{ if .Syntax.SupportsRequiredPrefix }}---|---|{{ end }}
{{ range $field := .Fields }}{{ fieldDoc . }}
{{ end }}

> JSON Demo

{{ jsonDemo . }}
{{ end }}
`
