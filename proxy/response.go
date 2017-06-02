package proxy

import "net/http"

type ResponseWriter struct {
	StatusCode int
	Output     []byte
	header     http.Header
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{}
}

func (rw *ResponseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = make(http.Header)
	}

	return rw.header
}

func (rw *ResponseWriter) Write(bytes []byte) (int, error) {
	if rw.StatusCode == 0 {
		rw.WriteHeader(200)
	}
	rw.Output = append(rw.Output, bytes...)

	return 0, nil
}

func (rw *ResponseWriter) WriteHeader(i int) {
	rw.StatusCode = i
}
