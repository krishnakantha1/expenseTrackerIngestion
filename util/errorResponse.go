package util

import "net/http"

func ForbidenResponse(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusForbidden)
}

func BadRequestResponse(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}
