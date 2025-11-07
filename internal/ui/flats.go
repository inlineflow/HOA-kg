package ui

import (
	"fmt"
	"hypermedia/internal/ui/pages"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

func (u *UI) FlatsCreate(w http.ResponseWriter, r *http.Request) {
	houseID, err := uuid.Parse(r.PathValue("house_id"))
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed parsing house_id from the URL: %v\n", err), 500)
		return
	}

	var opts []func(*templ.ComponentHandler)
	if r.Header.Get("HX-Request") != "" {
		opts = append(opts, templ.WithFragments("partial"))
	}

	templ.Handler(pages.CreateFlats(houseID), opts...).ServeHTTP(w, r)
}

func (u *UI) HandleCreateFlat(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	 fmt.Println(r.Form)
}
