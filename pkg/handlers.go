package eventgen

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
)

func (g *Generator) generateHandlers(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return fmt.Errorf("generateHandlers aborted due to context cancellation")
	default:
	}

	type d struct {
		DomainSchema
		Command    Command
		Package    string
		Emits      Event
		NewFile    bool
		ImportPath string
	}
	for _, domain := range g.app.Domains {

		for _, command := range domain.Commands {

			if command.Handler == "" {
				continue
			}

			handlerSplits := strings.Split(command.Handler, ":")

			var emits Event
			for _, event := range domain.Events {
				if event.Name == command.Emits {
					emits = event
				}
			}

			fileExists := false

			fileName := fmt.Sprintf("%s.go", strings.ToLower(handlerSplits[0]))

			funcName := command.Name
			if _, err := os.Stat(fileName); err == nil {
				fileExists = true

				foundFuncs, _ := parseFunctionsFromFile(fileName)
				if _, ok := foundFuncs[funcName]; ok {
					continue
				}

			}

			t := loadTemplate("handler_command")
			var buf bytes.Buffer
			if err := t.Execute(&buf, d{
				Package:      handlerSplits[1],
				Command:      command,
				DomainSchema: domain,
				Emits:        emits,
				ImportPath:   g.app.Package,
				NewFile:      !fileExists,
			}); err != nil {
				return err
			}

			fileBytes, err := formatAndImports(buf.Bytes())
			if err != nil {
				return err
			}

			createDirIfNotExists(fileName)

			if fileExists {
				f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}
				defer f.Close()
				if _, err := f.Write(fileBytes); err != nil {
					return err
				}
			} else {
				err = os.WriteFile(fileName, fileBytes, 0644)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}
