package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

type Template struct {
	// State
	raw string

	// Config
	Path string
}

// Template method that loads file
// The argument fsys is an abstraction
// NOTE: be sure to handle the error
func (t *Template) Load(fsys fs.FS) error {
	bytes, err := fs.ReadFile(fsys, t.Path)
	if err != nil {
		return err
	}

	t.raw = string(bytes)
	return nil
}

// Compile (Render and Save) template
// Load must be called first
// Can be called multiple times for one template with different data
func (t *Template) Compile(path string, data any) error {
	tmpl, err := template.New("").Parse(t.raw)
	if err != nil {
		return err
	}

	// Create directory containing path if it doesn't exist
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	// Create file
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	tmpl.Execute(file, data)
	return nil
}

// Convenience function to load all templates in a map (dict) of templates
// This is used instead of repeating t.Load() and error handling elsewhere
func LoadTemplates(templates map[string]*Template, fsys fs.FS) error {
	// Iterate over all templates
	for key, t := range templates {
		err := t.Load(fsys)
		if err != nil {
			return fmt.Errorf("Could not load template \"%s\" at path \"%s\" due to error \"%s\"", key, t.Path, err.Error())
		}
	}
	return nil
}
