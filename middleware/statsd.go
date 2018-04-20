package middleware

import (
	"net/http"
	"strings"

	"fmt"
	"github.com/Sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
	"regexp"
)

// Statsd is statsd metrics middleware.
type Statsd struct {
	logger  *logrus.Logger
	client  *statsd.Client
	withLog bool
}

// NewStatsd construct Statsd.
func NewStatsd(logger *logrus.Logger, client *statsd.Client, withLog bool) *Statsd {
	return &Statsd{
		logger:  logger,
		client:  client,
		withLog: withLog,
	}
}

// RegisterMetrics send metrics to statsd
func (s *Statsd) RegisterMetrics(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		path := req.URL.Path
		handlerName := prepareHandlerName(path)

		requestTimer := s.client.NewTiming()

		lrw := &LoggedResponseWriter{responseWriter: resp}

		handler.ServeHTTP(lrw, req)

		requestTimer.Send(fmt.Sprintf(
			"request.%v.%v.%v.request_time",
			req.Method,
			lrw.Status(),
			handlerName,
		))

		if s.withLog {
			//Example: 200 POST /rec/ (127.0.0.1) 1.460s
			s.logger.Infof("%v %v %v (%v) %.3fs",
				lrw.Status(), req.Method, path, req.RemoteAddr, requestTimer.Duration().Seconds())
		}
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
