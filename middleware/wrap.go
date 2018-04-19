package middleware

import (
	"net/http"
	"runtime/debug"
	"github.com/Sirupsen/logrus"
)

type Wrap struct {
	logger *logrus.Entry
	list   []func(http.Handler) http.Handler
	len    int8
}

func NewWrap(logger *logrus.Logger) *Wrap {
	return &Wrap{
		logger: logger.WithField("component", "middlewareWrap"),
		list:   make([]func(http.Handler) http.Handler, 0),
		len:    0,
	}
}

func (s *Wrap) Add(middleware func(http.Handler) http.Handler) {
	s.list = append(s.list, middleware)
	s.len++
}

func (s *Wrap) Do(handler http.Handler) http.Handler {
	for i := int8(0); i < s.len; i++ {
		handler = s.list[i](handler)
	}
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Errorf("Panic: %+v\n%s", err, debug.Stack())
			}
		}()

		handler.ServeHTTP(resp, req)
	})
}
