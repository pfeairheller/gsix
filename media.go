package gsix

import (
	"strings"
	"strconv"
	"sort"
)

type MediaRange struct {
	value string
	mtype string
	subtype string
	quality float32
	originalIndex int
	params map[string]string
}

type MediaRanges []*MediaRange

func (r MediaRanges) Len() int { return len(r) }
func (r MediaRanges) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r MediaRanges) Less(i, j int) bool {
	if r[i].quality == r[j].quality {
		return r[i].originalIndex < r[j].originalIndex
	} else {
		return r[i].quality < r[j].quality
	}
}

var (

	extnameMap = map[string]string {

    ".jpg":    "image/jpg"                              ,
    ".ico":    "image/vnd.microsoft.icon"               ,
    ".png":    "image/png"                              ,
    ".gif":    "image/gif"                              ,
    ".svg":    "image/svg"                              ,
    ".pdf":    "application/pdf"                        ,
    ".rb":     "application/ruby"                       ,
    ".txt":    "text/plain; charset=UTF-8"              ,
    ".css":    "text/css; charset=UTF-8"                ,
    ".js":     "application/javascript; charset=UTF-8"  ,
    ".coffee": "application/coffeescript; charset=UTF-8",
    ".sass":   "text/plain; charset=UTF-8"              ,
    ".scss":   "text/plain; charset=UTF-8"              ,
    ".html":   "text/html; charset=UTF-8"               ,
    ".woff":   "application/font-woff"                  ,
	}
)

func Extname(ext string) string {
	ll := strings.ToLower(ext)
	if !strings.HasPrefix(ll, ".") {
		ll = "." + ll
	}

	return extnameMap[ll]
}

func ParseAcceptHeader(header string) (MediaRanges) {
	var out MediaRanges

	for idx, mediaStr := range strings.Split(header, ",") {
		mr := parseAccept(mediaStr, idx)
		out = append(out, mr)
	}

	sort.Sort(out)
	return out
}

func (mediaRange *MediaRange) matches(media string) bool {
	return true
}

func parseAccept(mediaStr string, idx int) *MediaRange {
	mr := new(MediaRange)
	mr.originalIndex = idx
	p := strings.SplitN(mediaStr, ";", 2)
	types := strings.Split(p[0], "/")
	mr.mtype = types[0]
	if len(types) > 1 {
		mr.subtype = types[1]
	} else {
		mr.subtype = "*"
	}
	
	mr.value = p[0]
	if len(p) > 1 {
		paramString := strings.Trim(p[1], " ")
		params := strings.Split(paramString, ";")
		var acceptParams []string
		for _, param := range params {
			nameValue := strings.Split(param, "=")
			if len(nameValue) == 2 {
				if nameValue[0] == "q" {
					q, _ := strconv.ParseFloat(nameValue[1], 32)
					mr.quality = float32(q)
				} else {
					acceptParams = append(acceptParams, param)
				}
			}
		}
		
		if len(acceptParams) > 0 {
			mr.value = mr.value + ";" + strings.Join(acceptParams, ";")
		}
	} else {
		mr.quality = 1
	}

	return mr
}

func normalizeType(key string) *MediaRange {
	if strings.Contains(key, "/") {
		return parseAccept(key, -1)
 	} else {
		return &MediaRange{ extnameMap[key], "", "", 0, -1, nil }
	}
}

func normalizeTypes(keys []string) []*MediaRange {
	out := []*MediaRange{}
	for _, key := range keys {
		out = append(out, normalizeType(key))
	}

	return out
}
