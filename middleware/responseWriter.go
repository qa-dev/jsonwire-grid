package middleware

import (
	"net/http"
)

// LoggedResponseWriter - wraps http.ResponseWriter.
// Supports the ability to read a previously recorded response status.
type LoggedResponseWriter struct {
	responseWriter http.ResponseWriter
	status         int
}

// Header returns the result of calling the Header() method of the wrapped http.ResponseWriter.
func (lrw *LoggedResponseWriter) Header() http.Header {
	return lrw.responseWriter.Header()
}

// Write returns the result of calling the Write() method of the wrapped http.ResponseWriter.
func (lrw *LoggedResponseWriter) Write(data []byte) (int, error) {
	return lrw.responseWriter.Write(data)
}

// WriteHeader returns the result of calling the WriteHeader() method of the wrapped http.ResponseWriter
func (lrw *LoggedResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.responseWriter.WriteHeader(status)
}

// Status - returns the previously recorded response status.
func (lrw *LoggedResponseWriter) Status() int {
	return lrw.status
}

// ResponseWriter returns wrapped ResponseWriter.
func (lrw *LoggedResponseWriter) ResponseWriter() http.ResponseWriter {
	return lrw.responseWriter
}
