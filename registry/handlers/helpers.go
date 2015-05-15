package handlers

import (
	"encoding/json"
	"github.com/docker/distribution/registry/api/errcode"
	"io"
	"net/http"
)

// serveJSON marshals v and sets the content-type header to
// 'application/json'. If a different status code is required, call
// ResponseWriter.WriteHeader before this function.
func serveJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	sc := http.StatusInternalServerError

	if errs, ok := v.(errcode.Errors); ok && errs.Len() > 0 {
		sc = errs.Errors[0].Code.Descriptor().HTTPStatusCode
		if sc == 0 {
			sc = http.StatusInternalServerError
		}
	}

	w.WriteHeader(sc)

	enc := json.NewEncoder(w)

	if err := enc.Encode(v); err != nil {
		return err
	}

	return nil
}

// closeResources closes all the provided resources after running the target
// handler.
func closeResources(handler http.Handler, closers ...io.Closer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, closer := range closers {
			defer closer.Close()
		}
		handler.ServeHTTP(w, r)
	})
}
