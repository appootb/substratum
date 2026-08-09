package markdown

const methodTpl = `
<h3 id="{{ anchorName . }}">{{ headerTitle . }}</h3>

> {{ leadingComment .SourceCodeInfo }}

{{ $webDoc := (webDoc .) }}
{{ if $webDoc }}
* HTTP

	* URL: {{ $webDoc.URL }}
	* Method: {{ $webDoc.Method }}
{{ if $webDoc.ContentType }}	* Content-Type: {{ $webDoc.ContentType }}{{ end }}
{{ end }}
* Request Type: ***{{ inputMessage . }}***

{{ template "message" .Input }}

* Response Type: ***{{ outputMessage . }}***

{{ template "message" .Output }}
`
