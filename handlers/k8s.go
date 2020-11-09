package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Ok returns ok
func Ok(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	resp.WriteHeader(http.StatusOK)
}
