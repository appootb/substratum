package goshared

// boolTpl wraps const rules with proto2 optional presence checks.
const boolTpl = `{{ $opt := optional .Field }}
	{{ if $opt }}if m.{{ name .Field }} != nil { {{ end }}
	{{ template "const" . }}
	{{ if $opt }} } {{ end }}
`
