package goshared

const numTpl = `{{ $opt := optional .Field }}
	{{ if $opt }}if m.{{ name .Field }} != nil { {{ end }}
	{{ if .Rules.GetIgnoreEmpty }}
		if {{ accessor . }} != 0 {
	{{ end }}

	{{ template "const" . }}
	{{ template "ltgt" . }}
	{{ template "in" . }}

	{{ if .Rules.GetIgnoreEmpty }}
		}
	{{ end }}
	{{ if $opt }} } {{ end }}
`
