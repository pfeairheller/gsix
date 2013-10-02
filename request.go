package gsix

import (
	"net/http"
	_ "fmt"
	"errors"
	"strings"
)

type GRequest struct {
	App *GSix
	raw *http.Request
	accepts MediaRanges
	Url string
}

func NewGRequest(raw *http.Request, app *GSix) (* GRequest) {
	out := new(GRequest)
	out.raw = raw
	out.App = app
	out.Url = raw.URL.Path
	return out
}

func (req *GRequest) RawRequest() (*http.Request) {
	return req.raw
}

func (req *GRequest) Accepts(mtypes []string) (string, error) {
	normalizedTypes := normalizeTypes(mtypes)
	req.accepts = ParseAcceptHeader(req.raw.Header["Accept"][0])

	if len(req.accepts) == 0 {
		return mtypes[0], nil
	}
	
	for _, accept := range req.accepts {
		for jdx, mtype := range normalizedTypes {
			if req.Accept(mtype, accept) {
				return mtypes[jdx], nil
			}
		}
	}
	
	return "", errors.New("undefined")
}

func (req *GRequest) Accept(mtype *MediaRange, other *MediaRange) bool {
	t := strings.Split(mtype.value, "/")
	return (t[0] == other.mtype || "*" == other.mtype) && (t[1] == other.subtype || "*" == other.subtype)
	// && paramsEqual(mtype.params, other.params)
}

func (req *GRequest) Next(err error) bool {
	return true
}


