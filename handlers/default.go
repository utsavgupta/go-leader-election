package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Default handler returns Hello, World
func Default(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	text := []byte("Hello, World")

	w.WriteHeader(http.StatusOK)
	w.Write(text)
}
