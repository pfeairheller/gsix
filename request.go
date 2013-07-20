package gsix

import (
	"net/http"
)

type GRequest struct {
	raw *http.Request
}

func NewGRequest(raw *http.Request) (* GRequest) {
	out := new(GRequest)
	out.raw = raw

	return out
}

func (req *GRequest) RawRequest() (*http.Request) {
	return req.raw
}
