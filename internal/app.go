package eventgen

import (
	"bytes"
	"context"
	"fmt"
	"os"
)

type App struct {
	Package  string
	Domains  []DomainSchema
	EventMap map[string]string `yaml:"-"`
}

func (g *Generator) generateApp(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return fmt.Errorf("generateApp aborted due to context cancellation")
	default:
	}

	t := loadTemplate("app")

	var buf bytes.Buffer
	if err := t.Execute(&buf, g.app); err != nil {
		fmt.Println(buf.String())
		return err
	}

	fileBytes, err := formatAndImports(buf.Bytes())
	if err != nil {
		return err
	}

	err = os.WriteFile("gen/app.go", fileBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
