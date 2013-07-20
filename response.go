package gsix

import (
	"net/http"
)

type GResponse struct {
	raw http.ResponseWriter
}

func NewGResponse(raw http.ResponseWriter) (*GResponse){
	out := new(GResponse)
	out.raw = raw

	return out
}


func (resp *GResponse) RawWriter() (http.ResponseWriter) {
	return resp.raw
}

func (resp *GResponse) Status(status int) {
	resp.RawWriter().WriteHeader(status)
}

func (resp *GResponse) Set(name, value string) {
	resp.RawWriter().Header().Set(name, value)
}

func (resp *GResponse) Get(name string) (string){
	return resp.RawWriter().Header().Get(name)
}


