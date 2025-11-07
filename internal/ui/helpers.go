package ui

import (
	"github.com/a-h/templ"
	"net/http"
)

func RenderAndServe(w http.ResponseWriter, r *http.Request, t templ.Component) {

	var opts []func(*templ.ComponentHandler)
	if r.Header.Get("HX-Request") != "" {
		opts = append(opts, templ.WithFragments("partial"))
	}

	templ.Handler(t, opts...).ServeHTTP(w, r)
}
