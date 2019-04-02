package utilities

import "net/http"

type RequestHandler struct {
	handlers *map[string]func(http.ResponseWriter, *http.Request);
}

func NewRequestHandler(routers *map[string]func(http.ResponseWriter, *http.Request)) *RequestHandler {
	ret := new(RequestHandler);
	if ret == nil {
		return nil;
	}

	ret.handlers = routers;
	return ret;
}

func pathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		// should not happen
		return false
	}
	n := len(pattern)
	if pattern[n-1] != '/' {
		return pattern == path
	}
	return len(path) >= n && path[0:n] == pattern
}

func (self RequestHandler)findMacro(path string) func(w http.ResponseWriter, r *http.Request) {
	var n = 0
	var f func(http.ResponseWriter, *http.Request);
	for k, v := range *self.handlers {
		if !pathMatch(k, path) {
			continue
		}
		if f == nil || len(k) > n {
			n = len(k)
			f = v
		}
	}
	return f;
}

func (self RequestHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path;
	f, ok := (*self.handlers)[path];
	if ok {
		f(w, r);
		return;
	}

	f = self.findMacro(path);
	if f == nil {
		f = http.NotFound;
	}
	f(w, r);
}

