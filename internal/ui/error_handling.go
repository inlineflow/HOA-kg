package ui

import (
	"fmt"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error, code int) {
	fmt.Println(err)
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}
