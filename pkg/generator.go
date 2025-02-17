package eventgen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Field struct {
	Name   string  `yaml:"name"`
	Type   string  `yaml:"type"`
	GoType string  `yaml:"-"`
	Fields []Field `yaml:"fields"`
}

type Generator struct {
	app *App

	errorChan chan error
}

func New() *Generator {
	return &Generator{
		errorChan: make(chan error),
	}
}

func Generate() error {

	app, err := parse()
	if err != nil {
		return err
	}

	g := New()
	g.app = app

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	done := make(chan struct{})
	var retErr error

	// Handle any errors from goroutines
	go func() {
		select {
		case err := <-errChan:
			retErr = err
			cancel()
		case <-done:
		}
	}()

	wg.Add(7)
	go func() {
		defer wg.Done()
		if err := g.generateEntity(ctx); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := g.generateService(ctx); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := g.generateDomainEvents(ctx); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := g.generateReactors(ctx); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := g.generateEvent(ctx); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := g.generateApp(ctx); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := g.generateHandlers(ctx); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	if retErr != nil {
		return retErr
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
