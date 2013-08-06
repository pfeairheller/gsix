package gsix

import (
	"net/http"
)


type GSix struct {
	vars    map[string]string
	locals  map[string]string
	routes  map[string] *GHandler
	params  map[string] ParamCallback
	engines map[string] EngineCallback
	cache   map[string] *View
}

func NewGSix() (*GSix) {
	out := new(GSix)
	out.vars = make(map[string]string)
	out.routes = make(map[string]*GHandler)
	out.params = make(map[string]ParamCallback)
	out.engines = make(map[string]EngineCallback)
	out.cache = make(map[string]*View)

	return out
}


type EngineCallback func(path string, options map[string]string, callback ViewCallback)
type ParamCallback func(req *GRequest, resp GResponse, next func(), id string)(interface{})
type ViewCallback func(err error, html string) (bool)

func(g *GSix) CreateServer() (*Server) {
	server := NewServer()
	return server
}

func (g *GSix)Local(name string, value string) {
	g.locals[name] = value
}

func (g *GSix)Set(name string, value string) {
	g.vars[name] = value
}

func (g *GSix)Val(name string) (string){
	return g.vars[name]
}

func (g *GSix)Enable(name string) {
	g.vars[name] = "true"
}

func (g *GSix)Disable(name string) {
	delete(g.vars, name)
}

func (g *GSix)Enabled(name string) bool {
	return g.vars[name] == "true"
}

func (g *GSix)Disabled(name string) bool {
	return g.vars[name] != "true"
}


func(g *GSix) Map(handlerFunc GHandlerFunc, path string, method int) {
	if path == "" {
		path = "/"
	}

	handler := g.routes[path]
	if handler == nil {
		handler = NewGHandler(method)
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

func(g *GSix) Engine(ext string, engine EngineCallback) {
	if ext[0] != '.' {
		ext = "." + ext
	}
	g.engines[ext] = engine
}

func (g *GSix) Param(param string, callback ParamCallback) {
	g.params[param] = callback
}

func (g *GSix) Render(viewName string, options map[string]string, fn ViewCallback) {
	opts := make(map[string]string)
	merge(opts, options)

	//TODO:  Something about locals, not sure what it is

	var cache string
	var ok bool
	if cache, ok = opts["view cache"]; ok == false {
		cache, opts["view cache"] = "false", "false"
	}

	var view *View
	if cache == "true" {
		view = g.cache[viewName]
	}


	if view == nil {
		//TODO - CREATE VIEW HERE...
	}

	view.Render(opts, fn)
	
}

func merge(dst map[string]string, src map[string]string) {
	for k,v := range src {
		dst[k] = v
	}
}

