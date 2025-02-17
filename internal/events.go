package eventgen

import (
	"bytes"
	"context"
	"fmt"
	"os"
)

type Event struct {
	Name        string  `yaml:"name"`
	Type        string  `yaml:"type"`
	Handler     string  `yaml:"handler"`
	State       bool    `yaml:"state"`
	Description string  `yaml:"description"`
	Fields      []Field `yaml:"fields"`
}

func (g *Generator) generateEvent(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return fmt.Errorf("generateEvent aborted due to context cancellation")
	default:
	}

	t := loadTemplate("event_source")

	var buf bytes.Buffer
	if err := t.Execute(&buf, nil); err != nil {
		return err
	}

	fileBytes, err := formatAndImports(buf.Bytes())
	if err != nil {
		return err
	}

	err = os.WriteFile("gen/event_source.go", fileBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
