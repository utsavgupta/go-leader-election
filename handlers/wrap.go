package handlers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/utsavgupta/go-leader-election/globals"
	"github.com/utsavgupta/go-leader-election/logger"
)

// Wrap acts as a wrap around httprouter handlers which adds trace id to the incoming context
// and logs the time taken to serve a request
func Wrap(f httprouter.Handle) httprouter.Handle {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {

		start := time.Now()

		traceID := req.Header[logger.TRACEID]

		if traceID == nil {
			hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d%d", req.URL.String(), time.Now().UnixNano(), rand.Intn(500))))

			traceID = []string{fmt.Sprintf("%x", hash[:16])}
		}

		// set server name
		resp.Header().Set("Server", fmt.Sprintf("%s-%s", globals.APPLICATIONNAME, globals.APPLICATIONVERSION))

		// create a new context with the trace id
		ctx := context.WithValue(req.Context(), logger.KeyTraceID, traceID[0])

		// call the function handler with new request object
		f(resp, req.WithContext(ctx), params)

		// Log the latency of the request
		globals.Logger.Infof(ctx, "Request %s served in %d ms", req.URL.String(),
			time.Since(start)/time.Millisecond)
	}
}
