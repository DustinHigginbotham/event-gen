package eventgen

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
)

type ReactorSchema struct {
	Reactors []Reactor `yaml:"reactors"`
}

type Reactor struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	ReactsTo    string `yaml:"reactsTo"`
	Emits       string `yaml:"emits"`
	Type        string `yaml:"type"`
	Handler     string `yaml:"handler"`

	ActualEvent string `yaml:"-"`
}

func (g *Generator) generateReactors(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return fmt.Errorf("generateReactors aborted due to context cancellation")
	default:
	}

	type d struct {
		DomainSchema
		Package string
	}

	for _, domain := range g.app.Domains {

		if len(domain.Reactors) == 0 {
			continue
		}

		t := loadTemplate("reactor")
		var buf bytes.Buffer
		if err := t.Execute(&buf, d{Package: g.app.Package, DomainSchema: domain}); err != nil {
			fmt.Println(buf.String())
			return err
		}

		fileBytes, err := formatAndImports(buf.Bytes())
		if err != nil {
			return err
		}
		err = os.WriteFile(fmt.Sprintf("gen/%s.reactors.go", strings.ToLower(domain.Name)), fileBytes, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
