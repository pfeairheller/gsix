package gsix

import (
	"net/http"
)

type GResponse struct {
	raw http.ResponseWriter
	charset string
}

func NewGResponse(raw http.ResponseWriter) (*GResponse){
	out := new(GResponse)
	out.raw = raw
	out.charset = "utf-8"

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

func (resp *GResponse) Cookie(cookie *http.Cookie) {
	http.SetCookie(resp.RawWriter(), cookie)
}

func (resp *GResponse) Redirect(req *GRequest, urlStr string, code int) {
	http.Redirect(resp.RawWriter(), req.RawRequest(), urlStr, code)
}

func (resp *GResponse) Location(urlStr string) {
	resp.RawWriter().Header().Set("Location", urlStr)
	http.Error(resp.RawWriter(), "Found", 302)
}

func (resp *GResponse) Charset(value string) {
	resp.charset = value
}

func (resp *GResponse) Send(code int, body interface{}) {
	if resp.RawWriter().Header().Get("Content-Type") == "" {
		resp.RawWriter().Header().Set("Content-Type", "text/html; charset=" + resp.charset)
	}

	type := reflect.TypeOf(body)
	
}

