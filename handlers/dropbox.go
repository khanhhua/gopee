package handlers

import "net/http"

func Authorize(w http.ResponseWriter, r *http.Request) {
	if code, ok := r.URL.Query()["code"]; !ok {
		w.Write([]byte("Auth"))

		return
	}
}
