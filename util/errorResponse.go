package util

import (
	"net/http"
	"strings"
)

func ForbidenResponse(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusForbidden)
}

func BadRequestResponse(w http.ResponseWriter, msgs ...string) {
	var sb strings.Builder

	for _, v := range msgs {
		sb.WriteString(v)
		sb.WriteString(", ")
	}

	var msg string
	if sb.Len() > 0 {
		msg = sb.String()
		msg = msg[:len(msg)-2]
	} else {
		msg = "no issue message recoverd"
	}

	http.Error(w, msg, http.StatusBadRequest)
}
