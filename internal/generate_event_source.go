package eventgen

import (
	"bytes"
	"context"
	"fmt"
	"os"
)

// generateEventSource generates the event source file.
// This is the gen/event_source.go file.
// This includes the main event logic, which is not specific to any domain.
func (g *Generator) generateEventSource(ctx context.Context) error {

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
