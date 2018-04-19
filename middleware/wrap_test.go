package middleware

import (
	"testing"
	"github.com/Sirupsen/logrus"
	"net/http"
	"github.com/stretchr/testify/assert"
)

func TestWrap_Add(t *testing.T) {
	a := assert.New(t)
	wrap := NewWrap(logrus.StandardLogger())
	eLen := 5
	for i := 0; i < eLen; i++ {
		wrap.Add(func(handler http.Handler) http.Handler { return handler })
	}
	a.Len(wrap.list, eLen)
}

func TestWrap_Do(t *testing.T) {
	a := assert.New(t)
	wrap := NewWrap(logrus.StandardLogger())
	checkArray := make([]int, 0)

	wrap.Add(
		func(handler http.Handler) http.Handler {
			checkArray = append(checkArray, 1)
			return handler
		},
	)

	wrap.Add(
		func(handler http.Handler) http.Handler {
			checkArray = append(checkArray, 2)
			return handler
		},
	)

	wrap.Add(
		func(handler http.Handler) http.Handler {
			checkArray = append(checkArray, 3)
			return handler
		},
	)

	wrap.Do(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {}))

	a.Len(checkArray, 3)
	a.Equal(1, checkArray[0])
	a.Equal(2, checkArray[1])
	a.Equal(3, checkArray[2])
}
