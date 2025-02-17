package eventgen

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
)

func (g *Generator) generateService(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return fmt.Errorf("generateService aborted due to context cancellation")
	default:
	}

	type d struct {
		DomainSchema
		EventMap  map[string]string
		EntityMap map[string]string
		Package   string
	}

	for _, domain := range g.app.Domains {

		t := loadTemplate("service")

		entityMap := make(map[string]string)

		var buf bytes.Buffer
		if err := t.Execute(&buf, d{DomainSchema: domain, EventMap: g.app.EventMap, EntityMap: entityMap, Package: g.app.Package}); err != nil {
			return err
		}

		fileBytes, err := formatAndImports(buf.Bytes())
		if err != nil {
			return err
		}

		err = os.WriteFile(fmt.Sprintf("gen/%s.service.go", strings.ToLower(domain.Name)), fileBytes, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
