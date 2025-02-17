package eventgen

import (
	"embed"
	"fmt"
	"go/format"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

func loadTemplate(name string) *template.Template {
	data, err := templates.ReadFile(fmt.Sprintf("templates/%s.tmpl", name))
	if err != nil {
		panic(err)
	}
	return template.Must(template.New(name).Funcs(templateFunctions).Parse(string(data)))
}

func formatAndImports(src []byte) ([]byte, error) {
	formatted, err := format.Source(src)
	if err != nil {
		return nil, fmt.Errorf("failed to format code: %w", err)
	}
	return formatted, nil

	// Organize imports using `golang.org/x/tools/imports`
	// imported, err := imports.Process("", src, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to process imports: %w", err)
	// }

	// // Format the code using `go/format`
	// formatted, err := format.Source(imported)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to format code: %w", err)
	// }

	// return formatted, nil
}
