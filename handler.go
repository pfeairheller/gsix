package gsix

import (
	"net/http"
)

const (
	GET = 1 << iota
	HEAD
	POST
	PUT
	DELETE
	ALL = GET | HEAD | POST | PUT | DELETE
)

type GHandlerFunc func(*GRequest, *GResponse)

type GHandler struct {
	handlers []GHandlerFunc
	methods int
}

func NewGHandler(methods int) (*GHandler){
	handler := new(GHandler)
	handler.handlers = []GHandlerFunc{}
	handler.methods = methods
	return handler
}

func (handler *GHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	method := MethodConst(req.Method)
	if method & handler.methods > 0 {
		greq := NewGRequest(req)
		gresp := NewGResponse(resp)
		for _, handlerFunc := range handler.handlers {
			handlerFunc(greq, gresp)
		}
	}
}

func(handler *GHandler) Add(handlerFunc GHandlerFunc) {
	handler.handlers = append(handler.handlers, handlerFunc)
}

func MethodConst(method string) (int) {
	switch method {
	case "HEAD":
		return HEAD
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	default:
		return 0
	}
}
