package eventgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v2"
)

func parse() (*App, error) {
	var domains []DomainSchema

	eventMap := make(map[string]string)

	// Read the event-source directory
	files, err := os.ReadDir("event-source")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		data, err := os.ReadFile("event-source/" + file.Name())
		if err != nil {
			return nil, fmt.Errorf("error reading file from event-source: %s (%w)", file.Name(), err)
		}

		var domain DomainSchema
		err = yaml.Unmarshal(data, &domain)
		if err != nil {
			return nil, fmt.Errorf("error parsing file from event-source: %s (%w)", file.Name(), err)
		}

		if domain.Name == "" {
			domain.Name = strings.Title(strings.TrimSuffix(strings.TrimSuffix(file.Name(), ".yml"), ".yaml"))
		}

		for _, event := range domain.Events {
			eventMap[event.Type] = event.Name
		}

		for i, field := range domain.Entity.Fields {

			domain.Entity.Fields[i].GoType = field.Type
			if strings.Contains(field.Type, "array") {

				parts := strings.Split(field.Type, ";")
				if len(parts) == 2 {
					domain.Entity.Fields[i].GoType = parts[1]
					domain.Entity.Fields[i].Type = parts[0]
				}
			}
		}

		domains = append(domains, domain)
	}

	for _, domain := range domains {
		for i, reactor := range domain.Reactors {

			domain.Reactors[i].ActualEvent = eventMap[reactor.ReactsTo]

		}
	}

	// Read the go.mod file
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return nil, fmt.Errorf("error reading go.mod (%w)", err)
	}

	// Parse the go.mod file
	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("error parsing go.mod (%w)", err)
	}

	return &App{
		Package:  modFile.Module.Mod.Path,
		Domains:  domains,
		EventMap: eventMap,
	}, nil
}

func parseFunctionsFromFile(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, file, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	funcs := make(map[string]string)
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			name := fn.Name.Name
			params := []string{}
			for _, param := range fn.Type.Params.List {
				buf := &bytes.Buffer{}
				if err := format.Node(buf, fset, param.Type); err != nil {
					return nil, fmt.Errorf("failed to format parameter type: %w", err)
				}
				params = append(params, buf.String())
			}
			funcs[name] = strings.Join(params, ", ")
		}
	}

	return funcs, nil
}
