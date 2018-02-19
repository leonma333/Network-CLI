package network

import (
	"errors"
	"net/http"
)

// httpServer interface to define some http server functions
type httpServer interface {
	ListenAndServe(addr string, handler http.Handler) error
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// httpServerReal implements httpServer with real http calls
type httpServerReal struct{}

// httpServerMock implements httpServer with mock http calls
type httpServerMock struct {
	err bool
}

/*
 * Real http.ListenAndServe from httpServerReal (httpServer)
 */
func (h *httpServerReal) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

/*
 * Real http.Handle from httpServerReal (httpServer)
 */
func (h *httpServerReal) Handle(pattern string, handler http.Handler) {
	http.Handle(pattern, handler)
}

/*
 * Real http.HandleFunc from httpServerReal (httpServer)
 */
func (h *httpServerReal) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}

/*
 * Mock http.ListenAndServe from httpServerMock (httpServer)
 */
func (h *httpServerMock) ListenAndServe(addr string, handler http.Handler) error {
	// if err flag is set, return error
	if h.err {
		return errors.New("http.ListenAndServe error")
	}
	return nil
}

/*
 * Mock http.Handle from httpServerMock (httpServer)
 */
func (h *httpServerMock) Handle(pattern string, handler http.Handler) {
}

/*
 * Mock http.HandleFunc from httpServerMock (httpServer)
 */
func (h *httpServerMock) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
}
