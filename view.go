package gsix

import (
	"path/filepath"
	"os"
)

type View struct {
	name string
	root string
	ext string
	path string
	engine EngineCallback
}

func NewView(name, root string, engines *map[string] EngineCallback) (*View) {
	view := new(View)
	view.name = name
	view.ext = filepath.Ext(name)
	view.engine = (*engines)[view.ext]
	view.path = view.Lookup(name)

	return view
}

func (v *View) Lookup(name string) (string) {
	path := name
	if !filepath.IsAbs(path) {
		path = v.root + path
	}
	
	
	if exists(path) {
		return path
	} 

	path = filepath.Join(filepath.Dir(path), filepath.Base(path), "index", v.ext)

	return path
}

func (v *View) Render(options map[string]string, fn ViewCallback) {
	v.engine(v.path, options, fn)
}

func exists(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return false
}


