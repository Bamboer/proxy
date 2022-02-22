package router

import(
	"net/http"
	"net/http/httputil"
	"wikiproxy/pkg/common/log"
)


func Proxy(logger *log.Logger){
	proxy := httputil.ReverseProxy{
		Director: director
		ModifyResponse: modifyResponse
		FlushInterval: -1
		ErrorLog: logger
		ErrorHandler: errorHandler
	}
}

func director(req *http.Request){
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		req.Header.Set("User-Agent", "")
	}
}



func modifyResponse(res *http.Response) bool {
	if p.ModifyResponse == nil {
		return true
	}
	if err := p.ModifyResponse(res); err != nil {
		res.Body.Close()
		p.getErrorHandler()(rw, req, err)
		return false
	}
	return true
}