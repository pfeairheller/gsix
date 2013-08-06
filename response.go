package gsix

import (
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"strings"
)

type GResponse struct {
	raw http.ResponseWriter
	req *GRequest
	charset string
	status int
}

func NewGResponse(raw http.ResponseWriter, req *GRequest) (*GResponse){
	out := new(GResponse)
	out.raw = raw
	out.req = req
	out.charset = "utf-8"
	out.status = 200
	

	return out
}


func (resp *GResponse) RawWriter() (http.ResponseWriter) {
	return resp.raw
}

func (resp *GResponse) Status(status int) {
	resp.status = status
}

func (resp *GResponse) Set(name, value string) {
	resp.raw.Header().Set(name, value)
}

func (resp *GResponse) Get(name string) (string){
	return resp.raw.Header().Get(name)
}

func (resp *GResponse) Cookie(cookie *http.Cookie) {
	http.SetCookie(resp.raw, cookie)
}

func (resp *GResponse) Redirect(req *GRequest, urlStr string, code int) {
	http.Redirect(resp.raw, req.RawRequest(), urlStr, code)
}

func (resp *GResponse) Location(urlStr string) {
	resp.raw.Header().Set("Location", urlStr)
	http.Error(resp.raw, "Found", 302)
}

func (resp *GResponse) Charset(value string) {
	resp.charset = value
}

func (resp *GResponse) Type(value string) {
	if strings.Contains(value, "/") {
		resp.Set("Content-Type",  value)
	}	else if ext := Extname(value); ext != "" {
		resp.Set("Content-Type", ext)
	}
}

func (resp *GResponse) Send(code int, body string) {
	if resp.Get("Content-Type") == "" {
		resp.Set("Content-Type","text/html; charset=" + resp.charset)
	}

	resp.raw.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	resp.raw.WriteHeader(code)

	fmt.Fprintf(resp.raw, body)
}

func (resp *GResponse) SendStruct(code int, obj interface{}) {
	if resp.Get("Content-Type") == "" {
		resp.Set("Content-Type", "application/json")
	}

	body, err := json.Marshal(obj)
	if err != nil {
		log.Println("Error with marshalling ", err)
		resp.raw.WriteHeader(500)
		fmt.Fprintf(resp.raw, "Application Error")
		return
	}

	resp.raw.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	resp.raw.WriteHeader(code)

	resp.raw.Write(body)
}

type Format struct {
	mediaType string
	callback func()
}

func (resp *GResponse) Format(formats map[string]func()) {
	var fn, ok = formats["default"]
	delete(formats, "default")

	var keys []string
	for k,_ := range formats { keys = append(keys, k) }

	key, err := resp.req.Accepts(keys)

	if err == nil {
		resp.Set("Content-Type", normalizeType(key).value)
		formats[key]()
	} else if ok {
		fn()
	} else {
		//Send error
	}
	
}

func (resp *GResponse) Render(view string, locals map[string]string, callback ViewCallback) {
	
	var fn ViewCallback
	if callback == nil {
		fn = func(err error, html string) bool {
			if err != nil {
				return resp.req.Next(err)
			}

			resp.Send(200, html)
			return true
		}
	} else {
		fn = callback
	}

	resp.req.app.Render(view, locals, fn)
	
}


