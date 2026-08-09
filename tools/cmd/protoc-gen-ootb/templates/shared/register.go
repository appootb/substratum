package shared

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

func Register(tpl *template.Template, params pgs.Parameters) {
	fns := Funcs{
		Context: pgsgo.InitContext(params),
	}

	tpl.Funcs(map[string]interface{}{
		"package":          fns.PackageName,
		"serverName":       fns.ServerName,
		"serviceName":      fns.serviceName,
		"serverStream":     fns.ServerStream,
		"streamServerImpl": fns.streamServerImpl,
		"name":             fns.Name,
		"hasWebApi":        fns.hasWebApi,
		"hasWebStream":     fns.hasWebStream,
		"isWebStream":      fns.isWebStream,
		"webStreamURL":     fns.webStreamURL,
		"fullMethod":       fns.fullMethod,
		"subjectsName":     fns.subjectsName,
		"serviceSubjects":  fns.serviceSubjects,
		"rolesName":        fns.rolesName,
		"serviceRoles":     fns.serviceRoles,
		"serviceScope":     fns.serviceScope,
		"inputMessage":     fns.inputMessage,
		"outputMessage":    fns.outputMessage,
		"externalPackages": fns.externalPackages,
	})
}
