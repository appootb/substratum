package shared

import (
	"fmt"
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/appootb/substratum/v2/proto/go/api"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/genproto/googleapis/api/annotations"
)

const (
	TODOComment = "***TODO***"
)

type Funcs struct {
	pgsgo.Context
}

// headerTitle returns the documentation title for a node.
// Services use the proto service name (not the Go XxxServer interface name).
func (fns Funcs) headerTitle(node pgs.Node) string {
	if svc, ok := node.(pgs.Service); ok {
		return svc.Name().UpperCamelCase().String()
	}
	return fns.Name(node).String()
}

func (fns Funcs) anchorName(node pgs.Node) string {
	return strings.ToLower(fns.headerTitle(node))
}

func (fns Funcs) docFileName(file pgs.File) string {
	name := path.Base(file.Name().String())
	extName := path.Ext(name)
	return strings.ReplaceAll(name, extName, ".pb.md")
}

func (fns Funcs) tocComment(srcInfo pgs.SourceCodeInfo) string {
	comment := srcInfo.LeadingComments()
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return TODOComment
}

func (fns Funcs) leadingComment(srcInfo pgs.SourceCodeInfo) string {
	comment := TODOComment
	if srcInfo.LeadingComments() != "" {
		comment = srcInfo.LeadingComments()
	}
	comment = strings.Trim(comment, "\n")
	return strings.Replace(comment, "\n", "\n> ", -1)
}

func (fns Funcs) trailingComment(srcInfo pgs.SourceCodeInfo) string {
	comment := TODOComment
	if srcInfo.TrailingComments() != "" {
		comment = srcInfo.TrailingComments()
	} else if srcInfo.LeadingComments() != "" {
		comment = srcInfo.LeadingComments()
	}
	comment = strings.Trim(comment, "\n")
	return strings.Replace(comment, "\n", "<br>", -1)
}

func (fns Funcs) webURL(method pgs.Method) string {
	var (
		webStreamRule *api.WebsocketRule
		webApiRules   *annotations.HttpRule
	)

	if ok, _ := method.Extension(api.E_Websocket, &webStreamRule); ok {
		return webStreamRule.Url
	}
	if ok, _ := method.Extension(annotations.E_Http, &webApiRules); ok {
		switch p := webApiRules.Pattern.(type) {
		case *annotations.HttpRule_Get:
			return p.Get
		case *annotations.HttpRule_Put:
			return p.Put
		case *annotations.HttpRule_Post:
			return p.Post
		case *annotations.HttpRule_Delete:
			return p.Delete
		case *annotations.HttpRule_Patch:
			return p.Patch
		case *annotations.HttpRule_Custom:
			return p.Custom.Path
		}
	}
	return ""
}

type WebDoc struct {
	URL         string
	Method      string
	ContentType string
}

func (fns Funcs) webDoc(method pgs.Method) *WebDoc {
	var (
		webStreamRule *api.WebsocketRule
		webApiRules   *annotations.HttpRule
	)

	if ok, _ := method.Extension(api.E_Websocket, &webStreamRule); ok {
		return &WebDoc{
			URL:         fmt.Sprintf("`%s`", webStreamRule.Url),
			Method:      "`WS/WSS`",
			ContentType: "`TextFrame`",
		}
	}
	if ok, _ := method.Extension(annotations.E_Http, &webApiRules); ok {
		switch p := webApiRules.Pattern.(type) {
		case *annotations.HttpRule_Get:
			return &WebDoc{
				URL:    fmt.Sprintf("`%s`", p.Get),
				Method: fmt.Sprintf("`%s`", http.MethodGet),
			}
		case *annotations.HttpRule_Put:
			return &WebDoc{
				URL:         fmt.Sprintf("`%s`", p.Put),
				Method:      fmt.Sprintf("`%s`", http.MethodPut),
				ContentType: "`application/json`",
			}
		case *annotations.HttpRule_Post:
			return &WebDoc{
				URL:         fmt.Sprintf("`%s`", p.Post),
				Method:      fmt.Sprintf("`%s`", http.MethodPost),
				ContentType: "`application/json`",
			}
		case *annotations.HttpRule_Delete:
			return &WebDoc{
				URL:    fmt.Sprintf("`%s`", p.Delete),
				Method: fmt.Sprintf("`%s`", http.MethodDelete),
			}
		case *annotations.HttpRule_Patch:
			return &WebDoc{
				URL:         fmt.Sprintf("`%s`", p.Patch),
				Method:      fmt.Sprintf("`%s`", http.MethodPatch),
				ContentType: "`application/json`",
			}
		case *annotations.HttpRule_Custom:
			return &WebDoc{
				URL:         fmt.Sprintf("`%s`", p.Custom.Path),
				Method:      fmt.Sprintf("`%s`", p.Custom.Kind),
				ContentType: "`application/json`",
			}
		}
	}
	return nil
}

func (fns Funcs) enumName(enum pgs.Enum, pkg pgs.Package) string {
	if pkg.ProtoName() != enum.Package().ProtoName() &&
		fns.PackageName(pkg) != fns.PackageName(enum) {
		return fmt.Sprintf("%s.%s", fns.PackageName(enum).String(), fns.Name(enum).String())
	}
	return fns.Name(enum).String()
}

func (fns Funcs) messageName(message pgs.Message, pkg pgs.Package) string {
	if pkg.ProtoName() != message.Package().ProtoName() &&
		fns.PackageName(pkg) != fns.PackageName(message) {
		return fmt.Sprintf("%s.%s", fns.PackageName(message).String(), fns.Name(message).String())
	}
	return fns.Name(message).String()
}

func (fns Funcs) inputMessage(method pgs.Method) string {
	return fns.messageName(method.Input(), method.Package())
}

func (fns Funcs) outputMessage(method pgs.Method) string {
	return fns.messageName(method.Output(), method.Package())
}

func (fns Funcs) collectMessages(file pgs.File) []pgs.Message {
	seen := map[string]pgs.Message{}
	add := func(msgs ...pgs.Message) {
		for _, msg := range msgs {
			if msg == nil {
				continue
			}
			seen[msg.FullyQualifiedName()] = msg
			for _, dep := range msg.Dependencies() {
				seen[dep.FullyQualifiedName()] = dep
			}
		}
	}
	add(file.AllMessages()...)
	for _, svc := range file.Services() {
		for _, method := range svc.Methods() {
			add(method.Input(), method.Output())
		}
	}
	out := make([]pgs.Message, 0, len(seen))
	for _, msg := range seen {
		out = append(out, msg)
	}
	return out
}

func (fns Funcs) embedEnums(file pgs.File) []pgs.Enum {
	enums := map[string]pgs.Enum{}
	for _, enum := range file.Enums() {
		enums[enum.FullyQualifiedName()] = enum
	}
	for _, msg := range file.AllMessages() {
		for _, enum := range msg.Enums() {
			enums[enum.FullyQualifiedName()] = enum
		}
	}
	for _, msg := range fns.collectMessages(file) {
		for _, field := range msg.Fields() {
			t := field.Type()
			var enum pgs.Enum
			switch {
			case t.IsMap() && t.Element().IsEnum():
				enum = t.Element().Enum()
			case t.IsRepeated() && t.ProtoType() == pgs.EnumT:
				enum = t.Element().Enum()
			case t.ProtoType() == pgs.EnumT:
				enum = t.Enum()
			}
			if enum != nil {
				enums[enum.FullyQualifiedName()] = enum
			}
		}
	}
	resp := make([]pgs.Enum, 0, len(enums))
	for _, enum := range enums {
		resp = append(resp, enum)
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Name().String() < resp[j].Name().String()
	})
	return resp
}

func (fns Funcs) embedMessages(file pgs.File) []pgs.Message {
	resp := make([]pgs.Message, 0)
	for _, msg := range fns.collectMessages(file) {
		if msg.IsMapEntry() || msg.IsWellKnown() {
			continue
		}
		resp = append(resp, msg)
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Name().String() < resp[j].Name().String()
	})
	return resp
}
