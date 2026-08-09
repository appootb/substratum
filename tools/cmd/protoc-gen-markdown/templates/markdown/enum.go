package markdown

const enumTpl = `
<h3 id="{{ anchorName . }}">{{ headerTitle . }}</h3>

> {{ leadingComment .SourceCodeInfo }}

* Enum

|Name (string)|Value (integer)|Comment|
|---|---|---|
{{ range $v := .Values }}|{{ $v.Name }}|{{ $v.Value }}|{{ trailingComment $v.SourceCodeInfo }}|
{{ end }}
`
