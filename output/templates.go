package output

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/flyx/askew/data"
)

var fileHeader = template.Must(template.New("fileHeader").Funcs(template.FuncMap{
	"FormatImport": func(alias, path string) string {
		if filepath.Base(path) == alias {
			return "\"" + path + "\""
		}
		return alias + " \"" + path + "\""
	},
}).Parse(`
package {{.PackageName}}

import (
	"github.com/flyx/askew/runtime"
	"github.com/gopherjs/gopherjs/js"
	{{- range $alias, $path := .Imports }}
	{{FormatImport $alias $path}}{{ end }}
)
`))

var file = template.Must(template.New("file").Funcs(template.FuncMap{
	"Wrapper":      wrapperForType,
	"PathItems":    pathItems,
	"NameForBound": nameForBound,
	"Last":         last,
	"TWrapper": func(t *data.ParamType, name string) string {
		switch t.Kind {
		case data.IntType:
			return wrapperForType(data.IntVar) + "{BoundValue: " + name + "}"
		case data.StringType:
			return wrapperForType(data.StringVar) + "{BoundValue: " + name + "}"
		case data.BoolType:
			return wrapperForType(data.BoolVar) + "{BoundValue: " + name + "}"
		case data.PointerType:
			if t.ValueType.Kind == data.ObjectType {
				return wrapperForType(data.ObjectVar) + "{BoundValue: " + name + "}"
			}
		}
		panic("cannot gen wrapper for this type")
	},
	"IsBool": func(t *data.ParamType) bool {
		return t != nil && t.Kind == data.BoolType
	},
	"GenParams": func(params []data.Param) string {
		var items []string
		for _, p := range params {
			items = append(items, p.String())
		}
		return strings.Join(items, ", ")
	},
	"GenReturns": func(value *data.ParamType) string {
		if value == nil {
			return ""
		}
		return value.String()
	},
	"GenArgs": func(params []data.BoundParam) string {
		items := make([]string, 0, len(params))
		for _, p := range params {
			items = append(items, fmt.Sprintf("&p%s", p.Param))
		}
		return strings.Join(items, ", ")
	},
	"IsFormValue": func(bk data.BoundKind) bool {
		return bk == data.BoundFormValue
	},
	"IsEventValue": func(bk data.BoundKind) bool {
		return bk == data.BoundEventValue
	},
	"NeedsSelf": func(params []data.BoundParam) bool {
		for _, p := range params {
			if p.Value.Kind != data.BoundEventValue {
				return true
			}
		}
		return false
	},
	"TypeForKind": func(bk data.BoundKind) string {
		switch bk {
		case data.BoundProperty:
			return "BoundProperty"
		case data.BoundData:
			return "BoundData"
		case data.BoundClass:
			return "BoundClass"
		case data.BoundSelf:
			return "BoundSelf"
		default:
			panic("unknown BoundKind")
		}
	},
	"GenCallParams": func(params []data.Param) string {
		items := make([]string, 0, len(params))
		for _, p := range params {
			items = append(items, p.Name+" runtime.BoundValue")
		}
		return strings.Join(items, ", ")
	},
	"GenTypedArgs": func(params []data.Param) string {
		items := make([]string, 0, len(params))
		for _, p := range params {
			items = append(items, fmt.Sprintf("_%s.Get()", p.Name))
		}
		return strings.Join(items, ", ")
	},
	"GenComponentParams": func(params []data.ComponentParam) string {
		items := make([]string, 0, len(params))
		for _, p := range params {
			items = append(items, fmt.Sprintf("%s %s", p.Name, p.Type))
		}
		return strings.Join(items, ", ")
	},
	"ListParamVars": func(params []data.ComponentParam) string {
		items := make([]string, 0, len(params))
		for _, p := range params {
			items = append(items, p.Name)
		}
		return strings.Join(items, ", ")
	},
	"FieldType": func(e data.Embed) string {
		if e.T == "" {
			switch e.Kind {
			case data.OptionalEmbed:
				return "runtime.GenericOptional"
			case data.ListEmbed:
				return "runtime.GenericList"
			default:
				panic("unexpected field type")
			}
		}
		var b strings.Builder
		if e.Kind == data.DirectEmbed {
			b.WriteRune('*')
		}
		if e.Ns != "" {
			b.WriteString(e.Ns)
			b.WriteRune('.')
		}
		if e.Kind == data.OptionalEmbed {
			b.WriteString("Optional")
		}
		b.WriteString(e.T)
		if e.Kind == data.ListEmbed {
			b.WriteString("List")
		}
		return b.String()
	},
	"BlockNotEmpty": func(b data.Block) bool {
		return len(b.Assignments) > 0 || len(b.Controlled) > 0
	},
}).Option("missingkey=error").Parse(`
{{- define "Block"}}
  {{- range .Assignments}}
	{
		{{- if IsFormValue .Target.Kind}}
		var tmp runtime.BoundFormValue
		tmp.Init(runtime.WalkPath(block, {{PathItems .Path .Target.FormDepth}}), "{{.Target.ID}}", {{.Target.IsRadio}})
		{{- else}}
		var tmp runtime.{{TypeForKind .Target.Kind}}
		tmp.Init(runtime.WalkPath(block, {{PathItems .Path 0}}), "{{.Target.ID}}")
		{{- end}}
		runtime.Assign(&tmp, {{.Expression}})
	}
	{{- end}}

	{{- range .Controlled}}
	{{- if eq .Kind 0}}
	if {{.Expression}} {
		{{if BlockNotEmpty .Block}}
		block := runtime.WalkPath(block, {{PathItems .Path 0}})
		{{template "Block" .Block}}
		{{- end}}
	} else {
		_item := runtime.WalkPath(block, {{PathItems .Path 0}})
		_parent := _item.Get("parentNode")
		_parent.Call("replaceChild", js.Global.Get("document").Call("createComment", "removed"), _item)
	}
	{{- else }}
	{
		_orig := runtime.WalkPath(block, {{PathItems .Path 0}})
		_parent := _orig.Get("parentNode")
		_next := _orig.Get("nextSibling")
		_parent.Call("removeChild", _orig)
		for {{.Index}}{{with .Variable}}, {{.}}{{end}} := range {{.Expression}} {
			block := _orig.Call("cloneNode", true)
			{{template "Block" .Block}}
			_parent.Call("insertBefore", block, _next)
		}
	}
	{{- end}}
	{{- end}}
{{- end}}

{{- range .Components}}
{{- if .Controller}}
// {{.Name}}Controller can be implemented to handle external events
// generated by {{.Name}}
type {{.Name}}Controller interface {
	{{- range $name, $handler := .Controller }}
	{{$name}}({{GenParams $handler.Params }}){{GenReturns $handler.Returns}}
	{{- end }}
}
{{- end}}

// {{.Name}} is a DOM component autogenerated by Askew
type {{.Name}} struct {
	cd runtime.ComponentData
	{{- if .Controller }}
	// Controller is the adapter for events generated from this component.
	// if nil, events that would be passed to the controller will not be handled.
	Controller {{.Name}}Controller
	{{- end}}
	{{- range .Variables }}
	{{.Variable.Name}} {{Wrapper .Variable.Type}}
	{{- end}}
	{{- range $name, $type := .Fields}}
	{{$name}} {{$type}}
	{{- end}}
	{{- range .Embeds }}
	{{.Field}} {{FieldType .}}
	{{- end}}
}

// New{{.Name}} creates a new component and initializes it with Init.
func New{{.Name}}({{GenComponentParams .Parameters}}) *{{.Name}} {
	ret := new({{.Name}})
	ret.Init({{ListParamVars .Parameters}})
	return ret
}

// Data returns the object containing the component's DOM nodes.
// It implements the runtime.Component interface.
func (o *{{.Name}}) Data() *runtime.ComponentData {
	return &o.cd
}

// Init initializes the component, discarding all previous information.
// The component is initially a DocumentFragment until it gets inserted into
// the main document. It can be manipulated both before and after insertion.
func (o *{{.Name}}) Init({{GenComponentParams .Parameters}}) {
	o.cd.Init(runtime.InstantiateTemplateByID("{{.ID}}"))
	{{ range .Variables }}
	{{- if IsFormValue .Value.Kind}}
	o.{{.Variable.Name}}.BoundValue = runtime.NewBoundFormValue(&o.cd, "{{.Value.ID}}", {{.Value.IsRadio}}, {{PathItems .Path .Value.FormDepth}})
	{{- else}}
	o.{{.Variable.Name}}.BoundValue = runtime.New{{TypeForKind .Value.Kind}}(&o.cd, "{{.Value.ID}}", {{PathItems .Path 0}})
	{{- end}}
	{{- end}}
	{{- if BlockNotEmpty .Block}}
	{
		block := o.cd.Walk()
		{{- template "Block" .Block}}
	}
	{{- end}}
	{{- range .Embeds }}
	{
		container := o.cd.Walk({{PathItems .Path 1}})
		{{- if eq .Kind 0}}
		o.{{.Field}} = {{with .Ns}}{{.}}.{{end}}New{{.T}}({{.Args.Raw}})
		o.{{.Field}}.InsertInto(container, container.Get("childNodes").Index({{Last .Path}}))
		{{- if .Control}}
		o.{{.Field}}.Controller = o
		{{- end}}
		{{- else}}
		o.{{.Field}}.Init(container, {{Last .Path}})
		{{- if .Control}}
		o.{{.Field}}.DefaultController = o
		{{- end}}
		{{$e := .}}
		{{- range .ConstructorCalls}}
		{{- if eq .Kind 1}}
		if {{.Expression}} {
		{{- else if eq .Kind 2}}
		for {{.Index}}, {{.Variable}} := range {{.Expression}} {
		{{- end}}
		{{- if eq $e.Kind 2}}
		o.{{$e.Field}}.Set(
		{{- else}}
		o.{{$e.Field}}.Append(
		{{- end}}{{with $e.Ns}}{{.}}.{{end}}New{{$e.T}}({{.Args.Raw}}))
		{{- if ne .Kind 0}}
		}
		{{- end}}
		{{- end}}
		{{- end}}
	}
	{{- end}}
	{{- range .Captures}}
	{
		src := o.cd.Walk({{PathItems .Path 0}})
		{{- range .Mappings}}
		{
			wrapper := js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
				{{- if NeedsSelf .ParamMappings}}
				self := arguments[0].Get("currentTarget")
				{{- end}}
				{{- range .ParamMappings}}
				var p{{.Param}} runtime.{{NameForBound .Value.Kind}}
				{{- if IsFormValue .Value.Kind}}
				p{{.Param}}.Init(self.Call("closest", "form"), "{{.Value.ID}}", {{.Value.IsRadio}})
				{{- else if IsEventValue .Value.Kind}}
				p{{.Param}}.Init(arguments[0], "{{.Value.ID}}")
				{{- else}}
				p{{.Param}}.Init(self, "{{.Value.ID}}")
				{{- end}}
				{{- end}}
				{{- if eq .Handling 0}}
				o.call{{.Handler}}({{GenArgs .ParamMappings}})
				arguments[0].Call("preventDefault")
				{{- else if eq .Handling 2}}
				if o.call{{.Handler}}({{GenArgs .ParamMappings}}) {
					arguments[0].Call("preventDefault")
				}
				{{- else }}
				o.call{{.Handler}}({{GenArgs .ParamMappings}})
				{{- end}}
				return nil
			})
			src.Call("addEventListener", "{{.Event}}", wrapper)
		}
		{{- end}}
	}
	{{- end}}
	{{- if .Init}}
	o.init({{ListParamVars .Parameters}})
	{{- end}}
}

// InsertInto inserts this component into the given object.
// The component will be in inserted state afterwards.
//
// The component will be inserted in front of 'before', or at the end if 'before' is 'nil'.
func (o *{{.Name}}) InsertInto(parent *js.Object, before *js.Object) {
	o.cd.DoInsert(parent, before)
	{{- range .Embeds}}
	{{- if ne .Kind 0}}
	{{- if .T}}
	o.{{.Field}}.mgr.UpdateParent(o.cd.DocumentFragment(), parent, before)
	{{- else}}
	o.{{.Field}}.DoUpdateParent(o.cd.DocumentFragment(), parent, before)
	{{- end}}
	{{- end}}
	{{- end}}
}

// Extract removes this component from its current parent.
// The component will be in initial state afterwards.
func (o *{{.Name}}) Extract() {
	o.cd.DoExtract()
	{{- range .Embeds}}
	{{- if ne .Kind 0}}
	{{- if .T}}
	o.{{.Field}}.mgr.UpdateParent(o.cd.First().Get("parentNode"), o.cd.DocumentFragment(), nil)
	{{- else}}
	o.{{.Field}}.DoUpdateParent(o.cd.First().Get("parentNode"), o.cd.DocumentFragment(), nil)
	{{- end}}
	{{- end}}
	{{- end}}
}

{{$cName := .Name}}
{{- range $hName, $h := .Handlers}}
func (o *{{$cName}}) call{{$hName}}({{GenCallParams $h.Params}}) {{if IsBool $h.Returns}}bool{{end}} {
	{{- range $h.Params}}
	_{{.Name}} := {{TWrapper .Type .Name}}
	{{- end}}
	{{if IsBool $h.Returns}}return {{end}}o.{{$hName}}({{GenTypedArgs $h.Params}})
}
{{- end}}
{{- range $hName, $m := .Controller}}
{{- if $m.CanCapture}}
func (o *{{$cName}}) call{{$hName}}({{GenCallParams $m.Params}}) {{if IsBool $m.Returns}}bool{{end}} {
	if o.Controller == nil {
		return{{if IsBool $m.Returns}} false{{end}}
	}
	{{- range $m.Params}}
	_{{.Name}} := {{TWrapper .Type .Name}}
	{{- end}}
	{{if IsBool $m.Returns}}return {{end}}o.Controller.{{$hName}}({{GenTypedArgs $m.Params}})
}
{{- end}}
{{- end}}

{{if .NeedsList}}
// {{.Name}}List is a list of {{.Name}} whose manipulation methods auto-update
// the corresponding nodes in the document.
type {{.Name}}List struct {
	mgr runtime.ListManager
	items []*{{.Name}}
	{{- if .Controller}}
	DefaultController {{.Name}}Controller
	{{- end}}
}

// Init initializes the list, discarding previous data.
// The list's items will be placed in the given container, starting at the
// given index.
func (l *{{.Name}}List) Init(container *js.Object, index int) {
	l.mgr = runtime.CreateListManager(container, index)
	l.items = nil
}

// Len returns the number of items in the list.
func (l *{{.Name}}List) Len() int {
	return len(l.items)
}

// Item returns the item at the current index.
func (l *{{.Name}}List) Item(index int) *{{.Name}} {
	return l.items[index]
}

// Append appends the given item to the list.
func (l *{{.Name}}List) Append(item *{{.Name}}) {
	if item == nil {
		panic("cannot append nil to list")
	}
	l.mgr.Append(item)
	l.items = append(l.items, item)
	{{- if .Controller}}
	item.Controller = l.DefaultController
	{{- end}}
	return
}

// Insert inserts the given item at the given index into the list.
func (l *{{.Name}}List) Insert(index int, item *{{.Name}}) {
	var prev *js.Object
	if index < len(l.items) {
		prev = l.items[index].cd.First()
	}
	if item == nil {
		panic("cannot insert nil into list")
	}
	l.mgr.Insert(item, prev)
	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
	{{- if .Controller}}
	item.Controller = l.DefaultController
	{{- end}}
	return
}

// Remove removes the item at the given index from the list and returns it.
func (l *{{.Name}}List) Remove(index int) *{{.Name}} {
	item := l.items[index]
	l.mgr.Remove(item)
	copy(l.items[index:], l.items[index+1:])
	l.items = l.items[:len(l.items)-1]
	return item
}
{{end}}

{{- if .NeedsOptional}}
// Optional{{.Name}} is a nillable embeddable container for {{.Name}}.
type Optional{{.Name}} struct {
	cur *{{.Name}}
	mgr runtime.ListManager
	{{- if .Controller}}
	DefaultController {{.Name}}Controller
	{{- end}}
}

// Init initializes the container to be empty.
// The contained item, if any, will be placed in the given container at the
// given index.
func (o *Optional{{.Name}}) Init(container *js.Object, index int) {
	o.mgr = runtime.CreateListManager(container, index)
	o.cur = nil
}

// Item returns the current item, or nil if no item is assigned
func (o *Optional{{.Name}}) Item() *{{.Name}} {
	return o.cur
}

// Set sets the contained item removing the current one.
// Give nil as value to simply remove the current item.
func (o *Optional{{.Name}}) Set(value *{{.Name}}) {
	if o.cur != nil {
		o.mgr.Remove(o.cur)
	}
	o.cur = value
	if value != nil {
		o.mgr.Append(value)
		{{- if .Controller}}
		value.Controller = o.DefaultController
		{{- end}}
	}
}

{{- end}}
{{- end}}
`))

var skeleton = template.Must(template.New("skeleton").Funcs(template.FuncMap{
	"PathItems": pathItems,
	"Last":      last,
}).Parse(`
{{if .VarName}}
// {{.VarName}} holds the embedded components of the document's skeleton
var {{.VarName}} = struct {
	{{- range .Embeds}}
		// {{.Field}} is part of the main document.
		{{- if eq .Kind 0}}
			{{.Field}} *{{with .Ns}}{{.}}.{{end}}{{.T}}
		{{- else if eq .Kind 1}}
			{{- if .T}}
				{{.Field}} {{with .Ns}}{{.}}.{{end}}{{.T}}List
			{{- else}}
				{{.Field}} runtime.GenericList
			{{- end}}
		{{- else}}
			{{- if .T}}
				{{.Field}} {{with .Ns}}{{.}}.{{end}}Optional{{.T}}
			{{- else}}
				{{.Field}} runtime.GenericOptional
			{{- end}}
		{{- end}}
	{{- end -}}
}{
	{{- range .Embeds}}
		{{- if eq .Kind 0}}
			{{.Field}}: {{with .Ns}}{{.}}.{{end}}New{{.T}}({{.Args.Raw}}),
		{{- else if eq .Kind 1}}
			{{- if .T}}
				{{.Field}}: {{with .Ns}}{{.}}.{{end}}{{.T}}List{},
			{{- else}}
				{{.Field}}: runtime.GenericList{},
			{{- end}}
		{{- else}}
			{{- if .T}}
				{{.Field}}: {{with .Ns}}{{.}}.{{end}}Optional{{.T}}{},
			{{- else}}
				{{.Field}}: runtime.GenericOptional{},
			{{- end}}
		{{- end}}
	{{- end}}
}
{{- else}}
	{{range .Embeds}}
		// {{.Field}} is part of the main document.
		{{- if eq .Kind 0}}
			var {{.Field}} = {{with .Ns}}{{.}}.{{end}}New{{.T}}({{.Args.Raw}})
		{{- else if eq .Kind 1}}
			{{- if .T}}
				var {{.Field}} {{with .Ns}}{{.}}.{{end}}{{.T}}List
			{{- else}}
				var {{.Field}} runtime.GenericList
			{{- end}}
		{{- else}}
			{{- if .T}}
				var {{.Field}} {{with .Ns}}{{.}}.{{end}}Optional{{.T}}
			{{- else}}
				var {{.Field}} runtime.GenericOptional
			{{- end}}
		{{- end}}
	{{- end}}
{{- end}}

{{$varName := .VarName}}
func init() {
	html := js.Global.Get("document").Get("childNodes").Index(1)
	{{- range .Embeds}}
	{{- if eq .Kind 0}}
	{
		container := runtime.WalkPath(html, {{PathItems .Path 1}})
		{{with $varName}}{{.}}.{{end}}{{.Field}}.InsertInto(container, container.Get("childNodes").Index({{Last .Path}}))
	}
	{{- else}}
	{{with $varName}}{{.}}.{{end}}{{.Field}}.Init(runtime.WalkPath(html, {{PathItems .Path 1}}), {{Last .Path}})
	{{- end}}
	{{- end}}
}
`))
