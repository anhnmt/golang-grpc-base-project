package server

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}

func gatewayLoggerInterceptor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.Info()

		if r.RequestURI != "" {
			logger.Interface("method", r.RequestURI)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Err(err).Msg("Error reading body")
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		// Work / inspect body. You may even modify it!

		// And now set a new body, which will simulate the same data we read:
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Create a response wrapper:
		mrw := &MyResponseWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
		}

		logger.Interface("header", r.Header.Clone())

		if len(body) > 0 {
			logger.RawJSON("body", body)
		}

		h.ServeHTTP(mrw, r)

		logger.RawJSON("response", mrw.buf.Bytes())

		// Now inspect response, and finally send it out:
		// (You can also modify it before sending it out!)
		if _, err = io.Copy(w, mrw.buf); err != nil {
			log.Printf("Failed to send out response: %v", err)
		}

		logger.Msg("Log payload interceptor")
	})
}
