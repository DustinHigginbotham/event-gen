package eventgen

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func (g *Generator) generateEventHandlers() error {

	type d struct {
		DomainSchema
		Event      Event
		Entity     Entity
		Package    string
		NewFile    bool
		ImportPath string
	}

	for _, domain := range g.app.Domains {

		for _, ev := range domain.Events {

			if ev.Handler == "" {
				continue
			}

			handlerSplits := strings.Split(ev.Handler, ":")

			fileExists := false

			fileName := fmt.Sprintf("%s.go", strings.ToLower(handlerSplits[0]))

			funcName := fmt.Sprintf("Handle%s", strings.Title(ev.Name))
			if _, err := os.Stat(fileName); err == nil {
				fileExists = true

				foundFuncs, _ := parseFunctionsFromFile(fileName)
				if _, ok := foundFuncs[funcName]; ok {
					continue
				}

			}

			t := loadTemplate("handler_event")
			var buf bytes.Buffer
			if err := t.Execute(&buf, d{
				Package:      handlerSplits[1],
				Event:        ev,
				Entity:       domain.Entity,
				DomainSchema: domain,
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
