package ui

import "net/http"

func HandleError(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}
