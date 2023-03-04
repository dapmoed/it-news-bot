package template

import (
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	ErrNotFound = errors.New("not found template")
)

type Templates struct {
	files map[string]*template.Template
}

func New(path string) (*Templates, error) {
	result := &Templates{
		files: make(map[string]*template.Template),
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			tmpl, err := template.ParseFiles(filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}
			result.files[fileNameWithoutExt(filepath.Base(file.Name()))] = tmpl
		}
	}
	return result, nil
}

func fileNameWithoutExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func (t *Templates) Get(s string) (*template.Template, error) {
	if v, ok := t.files[s]; ok {
		return v, nil
	}
	return nil, ErrNotFound
}

func (t *Templates) Execute(name string, wr io.Writer, data any) error {
	tmpl, err := t.Get(name)
	if err != nil {
		return err
	}

	err = tmpl.Execute(wr, data)
	if err != nil {
		return err
	}
	return nil
}
