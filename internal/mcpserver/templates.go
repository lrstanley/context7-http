// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"embed"
	"maps"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	//go:embed templates
	templateDir embed.FS

	templates = template.Must(
		template.New("base").
			Funcs(sprig.FuncMap()).
			ParseFS(templateDir, "templates/*.gotmpl"),
	)
)

func (s *Server) render(name string, data map[string]any) (string, error) {
	var out strings.Builder

	merged := maps.Clone(s.baseVariables)
	maps.Copy(merged, data)

	err := templates.ExecuteTemplate(&out, name+".gotmpl", merged)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (s *Server) mustRender(name string, data map[string]any) string {
	out, err := s.render(name, data)
	if err != nil {
		panic(err)
	}
	return out
}
