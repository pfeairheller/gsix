package gsix

import (
	"net/http"
	_ "fmt"
	"errors"
)

type GRequest struct {
	raw *http.Request
	accepts MediaRanges
}

func NewGRequest(raw *http.Request) (* GRequest) {
	out := new(GRequest)
	out.raw = raw

	return out
}

func (req *GRequest) RawRequest() (*http.Request) {
	return req.raw
}

func (req *GRequest) Accepts(mtypes []string) (string, error) {
	if len(req.accepts) == 0 {
		return mtypes[0], nil
	}
	
	req.accepts = ParseAcceptHeader(req.raw.Header["Accept"][0])

	for idx, accept := range req.accepts {
		for jdx, mtype := range mtypes {
			
		}
	}
	
	
	return "", errors.New("undefined")
}

func (req *GRequest) Accept(type string) {

}


