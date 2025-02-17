package eventgen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type Field struct {
	Name   string  `yaml:"name"`
	Type   string  `yaml:"type"`
	GoType string  `yaml:"-"`
	Fields []Field `yaml:"fields"`
}

type Generator struct {
	app *App
}

func New() *Generator {
	return &Generator{}
}

func Generate() error {

	if err := createDirIfNotExists("./gen/"); err != nil {
		return err
	}

	app, err := parse()
	if err != nil {
		return err
	}

	g := New()
	g.app = app

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grp, ctx := errgroup.WithContext(ctx)

	generators := []func(ctx context.Context) error{
		g.generateEntity,
		g.generateService,
		g.generateDomainEvents,
		g.generateReactors,
		g.generateEvent,
		g.generateApp,
		g.generateHandlers,
	}

	for _, generator := range generators {
		grp.Go(func() error { return generator(ctx) })
	}

	if err := grp.Wait(); err != nil {
		return err
	}

	return nil
}

func createDirIfNotExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}
