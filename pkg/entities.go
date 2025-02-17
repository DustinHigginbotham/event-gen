package eventgen

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
)

type Entity struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Fields      []Field `yaml:"fields"`
}

func (g *Generator) generateEntity(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return fmt.Errorf("generateEntity aborted due to context cancellation")
	default:
	}

	type d struct {
		DomainSchema
		Package string
	}
	for _, domain := range g.app.Domains {

		t := loadTemplate("entity")
		var buf bytes.Buffer
		if err := t.Execute(&buf, d{Package: g.app.Package, DomainSchema: domain}); err != nil {
			return err
		}

		fileBytes, err := formatAndImports(buf.Bytes())
		if err != nil {
			return err
		}

		err = os.WriteFile(fmt.Sprintf("gen/%s.entity.go", strings.ToLower(domain.Name)), fileBytes, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
