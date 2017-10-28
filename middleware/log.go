package middleware

import (
	"net/http"
	"runtime/debug"
	"strings"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
	"regexp"
)

// LogMiddleware - wraps adds logging to the query handlers.
type LogMiddleware struct {
	statsd *statsd.Client
}

// NewLogMiddleware - constructor of LogMiddleware.
func NewLogMiddleware(statsd *statsd.Client) *LogMiddleware {
	return &LogMiddleware{
		statsd: statsd,
	}
}

// Log - wraps http.Handler for runtime logging.
func (m *LogMiddleware) Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				log.Fatalf("Panic: %+v\n%s", err, debug.Stack())
			}
		}()

		path := req.URL.Path
		handlerName := prepareHandlerName(path)

		requestTimer := m.statsd.NewTiming()

		lrw := &LoggedResponseWriter{responseWriter: resp}

		handler.ServeHTTP(lrw, req)

		requestTimer.Send(fmt.Sprintf(
			"request.%v.%v.%v.request_time",
			req.Method,
			lrw.Status(),
			handlerName,
		))

		//Example: 200 POST /rec/ (127.0.0.1) 1.460s
		log.Infof("%v %v %v (%v) %.3fs",
			lrw.Status(), req.Method, path, req.RemoteAddr, requestTimer.Duration().Seconds())

	})
}

func prepareHandlerName(name string) string {
	// todo: часть одного метода проебываем, но пока забьем на это
	r := regexp.MustCompile("(/session|/element|/window|/cookie|/attribute|/equals|/css|/key)(/[^/]*)")
	name = r.ReplaceAllString(name, "$1")
	name = strings.Trim(name, "/")
	name = strings.Replace(name, "/", "_", -1)
	return name

}
