// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"embed"
	"maps"
	"strings"
	"text/template"
)

var (
	//go:embed templates
	templateDir embed.FS

	templates = template.Must(
		template.New("base").
			Funcs(template.FuncMap{
				"sub": func(a, b int) int {
					return a - b
				},
			}).
			ParseFS(templateDir, "templates/*.gotmpl"),
	)
)

func (s *Server) render(name string, data any) (string, error) {
	var out strings.Builder

	if m, ok := data.(map[string]any); ok {
		merged := maps.Clone(s.baseVariables)
		maps.Copy(merged, m)
		data = merged
	}

	err := templates.ExecuteTemplate(&out, name+".gotmpl", data)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func (s *Server) mustRender(name string, data map[string]any) string {
	out, err := s.render(name, data)
	if err != nil {
		panic(err)
	}
	return out
}
