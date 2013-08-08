package gsix

import (
	"path/filepath"
	"os"
	"errors"
	"strings"
)

type View struct {
	name string
	root string
	ext string
	path string
	engine Engine
	defaultEngine string
}

func NewView(name, root, defaultEngine string, engines *map[string]Engine) (*View, error) {
	view := new(View)
	view.name = name
	view.root = root
	view.ext = filepath.Ext(name)
	view.defaultEngine = defaultEngine
	if view.defaultEngine == "" && view.ext == "" {
		return nil, errors.New("No engine or default engine")
	}

	if view.ext == "" {
		view.ext = "." + strings.TrimPrefix(view.defaultEngine, ".")
		view.name = view.name 
	}

	var ok bool
	view.engine, ok = (*engines)[view.ext]
	if !ok {
		return nil, errors.New("No suitable engine found for " + name)
	}

	view.path = view.Lookup(name)
	if view.path == "" {
		
		return nil, errors.New("Failed to lookup view " + name)
	}

	return view, nil
}

func (v *View) Lookup(name string) (string) {
	path := name
	if !filepath.IsAbs(path) {
		path = filepath.Join(v.root, path)
	}
	
	if exists(path) {
		return path
	} 

	path = filepath.Join(filepath.Dir(path), filepath.Base(path), "index", v.ext)

	return path
}

func (v *View) Render(data interface{}, options map[string]string, fn ViewCallback) (bool) {
	v.engine.Render(v.path, data, options, fn)
	return true
}

func exists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return false
}


