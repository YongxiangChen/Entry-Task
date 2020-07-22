package server

import "net/http"

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		index(w, r)
		return
	case "/login":
		login(w, r)
		return
	case "/userhome":
		userhome(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}
}
