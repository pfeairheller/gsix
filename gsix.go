package gsix

import (
	"net/http"
	"path/filepath"
	"os"
)


type GSix struct {
	vars    map[string] interface{}
	locals  map[string]string
	routes  map[string] *GHandler
	params  map[string] ParamCallback
	engines map[string] Engine
	cache   map[string] *View
}

func NewGSix() (*GSix) {
	out := new(GSix)
	out.vars = make(map[string]interface{})
	out.routes = make(map[string]*GHandler)
	out.params = make(map[string]ParamCallback)
	out.engines = make(map[string]Engine)
	out.cache = make(map[string]*View)

	//TODO  - default Configuration...

	dir, _ := os.Getwd()
	out.Engine("html", NewTemplateEngine())
	out.Set("views", filepath.Join(dir, "views"))
	out.Set("view engine", "html")

	return out
}


type ParamCallback func(req *GRequest, resp GResponse, next func(), id string)(interface{})
type ViewCallback func(err error, html string) (bool)


type Engine interface {
	Render (path string, data interface{}, options map[string]string, callback ViewCallback)
}

func(g *GSix) CreateServer() (*Server) {
	server := NewServer()

	return server
}

func (g *GSix)Local(name string, value string) {
	g.locals[name] = value
}

func (g *GSix)Set(name string, value interface{}) {
	g.vars[name] = value
}

func (g *GSix)Val(name string) (interface{}){
	return g.vars[name]
}

func (g *GSix)Enable(name string) {
	g.vars[name] = true
}

func (g *GSix)Disable(name string) {
	delete(g.vars, name)
}

func (g *GSix)Enabled(name string) bool {
	return g.vars[name] == true
}

func (g *GSix)Disabled(name string) bool {
	return g.vars[name] != true
}


func(g *GSix) Map(handlerFunc GHandlerFunc, path string, method int) {
	if path == "" {
		path = "/"
	}

	handler := g.routes[path]
	if handler == nil {
		handler = NewGHandler(method, g)
		g.routes[path] = handler
		http.Handle(path, handler)
	}

	handler.Add(handlerFunc)
}

func(g *GSix) Use(handlerFunc GHandlerFunc, path string) {
	g.Map(handlerFunc, path, ALL)
}

func(g *GSix) All(handlerFunc GHandlerFunc, path string) {
	g.Map(handlerFunc, path, ALL)
}

func(g *GSix) Get(handlerFunc GHandlerFunc, path string) {
	g.Map(handlerFunc, path, GET)
}

func(g *GSix) Post(handlerFunc GHandlerFunc, path string) {
	g.Map(handlerFunc, path, POST)
}

func(g *GSix) Put(handlerFunc GHandlerFunc, path string) {
	g.Map(handlerFunc, path, PUT)
}

func(g *GSix) Delete(handlerFunc GHandlerFunc, path string) {
	g.Map(handlerFunc, path, PUT)
}

func (g *GSix) Static(pathname string) (GHandlerFunc) {
	handler := http.FileServer(http.Dir(pathname))
	return func(req *GRequest, resp *GResponse) {
		handler.ServeHTTP(resp.RawWriter(), req.RawRequest())
	}
}

func(g *GSix) Engine(ext string, engine Engine) {
	if ext[0] != '.' {
		ext = "." + ext
	}
	g.engines[ext] = engine
}

func (g *GSix) Param(param string, callback ParamCallback) {
	g.params[param] = callback
}

func (g *GSix) Render(viewName string, data interface{}, options map[string]string, fn ViewCallback) (bool) {
	opts := make(map[string]string)
	merge(opts, options)

	//TODO:  Something about locals, not sure what it is

	var cache string
	var ok bool
	var err error
	if cache, ok = opts["view cache"]; ok == false {
		cache, opts["view cache"] = "false", "false"
	}

	var view *View
	if cache == "true" {
		view = g.cache[viewName]
	}


	if view == nil {
		view, err = NewView(viewName, g.Val("views").(string), g.Val("view engine").(string), &g.engines)
		if err != nil {
			return fn(err, "")
		}

		if cache == "true" {
			g.cache[viewName] = view
		}
		
	}

	return view.Render(data, opts, fn)
	
}

func merge(dst map[string]string, src map[string]string) {
	for k,v := range src {
		dst[k] = v
	}
}

