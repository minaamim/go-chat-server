package handler

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {

	_, _ = w.Write([]byte("ok"))

}
