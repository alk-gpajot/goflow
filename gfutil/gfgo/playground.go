package gfgo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"

	"github.com/alkemics/goflow"
)

func GeneratePlayground(w io.Writer, pkgPath string, nodes []Node) error {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName,
	}, pkgPath)
	if err != nil {
		return err
	}

	var pkg string
	for _, p := range pkgs {
		if p.PkgPath == pkgPath {
			pkg = p.Name
		}
	}

	if pkg == "" {
		return errors.New("could not extract package name")
	}

	ns := make([]Node, 0, len(nodes))
	imports := make([]goflow.Import, 0)
	for _, n := range nodes {
		if n.Method != "Run" {
			continue
		}

		imports = append(imports, n.Imports...)
		ns = append(ns, n)
	}

	sort.SliceStable(ns, func(i, j int) bool {
		lhf := fmt.Sprintf("%s.%s", ns[i].Pkg, ns[i].Typ)
		rhf := fmt.Sprintf("%s.%s", ns[j].Pkg, ns[j].Typ)
		return lhf <= rhf
	})

	imports = append(
		imports,
		goflow.Import{Pkg: "errors", Dir: "errors"},
		goflow.Import{Pkg: "fmt", Dir: "fmt"},
		goflow.Import{Pkg: "http", Dir: "net/http"},
		goflow.Import{Pkg: "json", Dir: "encoding/json"},
		goflow.Import{Pkg: "strings", Dir: "strings"},
	)

	dependencyMap := make(map[string]string)
	for _, node := range nodes {
		for _, dep := range node.Dependencies {
			dependencyMap[dep.Name] = dep.Type
		}
	}

	dependencies := make([]goflow.Field, 0, len(dependencyMap))
	for name, typ := range dependencyMap {
		dep := goflow.Field{
			Name: name,
			Type: typ,
		}
		dependencies = append(dependencies, dep)
	}

	sort.SliceStable(dependencies, func(i, j int) bool {
		return dependencies[i].Name <= dependencies[j].Name
	})

	tmplData := struct {
		Pkg          string
		PkgPath      string
		Imports      []goflow.Import
		Dependencies []goflow.Field
		Nodes        []Node
	}{
		Pkg:          pkg,
		PkgPath:      pkgPath,
		Imports:      imports,
		Dependencies: dependencies,
		Nodes:        ns,
	}

	buf := bytes.Buffer{}
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

var tmpl = template.Must(
	template.
		New("playground").
		Funcs(template.FuncMap{
			"JSONTag": func(s string) string {
				return fmt.Sprintf("`json:\"%s\"`", s)
			},
			"NameAndTypes": func(fields []goflow.Field) string {
				s := make([]string, len(fields))
				for i, f := range fields {
					s[i] = fmt.Sprintf("%s %s", f.Name, f.Type)
				}
				return strings.Join(s, ", ")
			},
			"Names": func(fields []goflow.Field) string {
				s := make([]string, len(fields))
				for i, f := range fields {
					s[i] = f.Name
				}
				return strings.Join(s, ", ")
			},
			"Public": func(s string) string {
				if s == "" {
					return ""
				}

				if len(s) == 1 {
					return strings.ToUpper(s)
				}

				return strings.ToUpper(string(s[0])) + s[1:]
			},
		}).
		Parse(`
// Code generated by goflow DO NOT EDIT.

// +build !codeanalysis

{{ $pkg := .Pkg -}}

package {{ .Pkg }}

import (
	{{ range .Imports -}}
	{{ .Pkg }} "{{ .Dir }}"
	{{ end -}}
)

type Playground struct{
	{{ range .Dependencies -}}
	{{ .Name }} {{ .Type }}
	{{ end -}}
}

func NewPlayground({{ NameAndTypes .Dependencies }}) Playground {
	return Playground{
        {{ range .Dependencies -}}
        {{ .Name }}: {{ .Name }},
        {{ end -}}
    }
}

func (p Playground) Run(ctx context.Context, name string, jsonInputs *json.RawMessage) (jsonOutputs *json.RawMessage, err error) {
	nodes := map[string]func(context.Context, *json.RawMessage) (*json.RawMessage, error) {
		{{ range .Nodes -}}
		"{{ .PkgPath }}.{{ .Typ }}": func(ctx context.Context, in *json.RawMessage) (*json.RawMessage, error) {
			if in == nil {
				return nil, nil
			}

			var params struct{
				{{ range .Inputs -}}
				{{ if ne .Type "context.Context" -}}
				{{ Public .Name }} {{ .Type }} {{ JSONTag .Name }}
				{{ end -}}
				{{ end -}}
			}
			if err := json.Unmarshal(*in, &params); err != nil {
				return nil, err
			}

			{{ if eq .Pkg $pkg -}}
			g := New{{ .Typ }}({{ range .Dependencies }}p.{{ .Name }},{{ end }})
			{{ else -}}
			g := {{ .Pkg }}.New{{ .Typ }}({{ range .Dependencies }}p.{{ .Name }},{{ end }})
			{{ end -}}
			{{ Names .Outputs }} := g.Run(
				{{ range .Inputs -}}
				{{ if eq .Type "context.Context" -}}
				ctx,
				{{ else -}}
				params.{{ Public .Name }},
				{{ end -}}
				{{ end -}}
			)

			res := map[string]interface{}{
				{{ range .Outputs -}}
				"{{ .Name }}": {{ .Name }},
				{{ end -}}
			}

			var out json.RawMessage
			if b, err := json.Marshal(res); err != nil {
				return nil, err
			} else {
				out = json.RawMessage(b)
			}

			return &out, nil
		},
		{{ end -}}
	}

	run := nodes[name]
	if run == nil {
		return nil, errors.New("not found")
	}
	return run(ctx, jsonInputs)
}

func (p Playground) List() (map[string]interface{}) {
	return map[string]interface{}{
		{{ range .Nodes -}}
		"{{ .PkgPath }}.{{ .Typ }}":
			{{- if eq .Pkg $pkg -}}
			New{{ .Typ }}({{ range .Dependencies }}p.{{ .Name }},{{ end }}),
			{{ else -}}
			{{ .Pkg }}.New{{ .Typ }}({{ range .Dependencies }}p.{{ .Name }},{{ end }}),
			{{ end -}}
		{{ end -}}
	}
}

func (p Playground) Get(name string) (interface{}) {
	nodes := p.List()
	return nodes[name]
}
`))
