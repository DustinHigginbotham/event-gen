package eventgen

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
)

// generateDomainEvents generates the events for the given domain.
// This is the gen/<domain>.events.go file.
func (g *Generator) generateDomainEvents(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return fmt.Errorf("generateDomainEvents aborted due to context cancellation")
	default:
	}

	type d struct {
		DomainSchema
		Package string
	}

	for _, domain := range g.app.Domains {

		t := loadTemplate("domain.event")
		var buf bytes.Buffer
		if err := t.Execute(&buf, d{Package: g.app.Package, DomainSchema: domain}); err != nil {
			return err
		}

		fileBytes, err := formatAndImports(buf.Bytes())
		if err != nil {
			return err
		}
		err = os.WriteFile(fmt.Sprintf("gen/%s.events.go", strings.ToLower(domain.Name)), fileBytes, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
