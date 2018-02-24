package network

import (
	"errors"
	"io"
	"net/http"
)

// httpServer interface to define some http server functions
type httpServer interface {
	ListenAndServe(addr string, handler http.Handler) error
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Get(url string) (resp *http.Response, err error)
}

// httpServerReal implements httpServer with real http calls
type httpServerReal struct{}

// httpServerMock implements httpServer with mock http calls
type httpServerMock struct {
	err  bool
	body io.ReadCloser
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
 * Real http.Get from httpServerReal (httpServer)
 */
func (h *httpServerReal) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
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

// Mock http.Get from httpServerMock (httpServer)
func (h *httpServerMock) Get(url string) (resp *http.Response, err error) {
	if h.err {
		return nil, errors.New("http.Get error")
	}
	return &http.Response{
		Body: h.body,
	}, nil
}

// readerMock implements ReadCloser with mock Reader
type readerMock struct {
	io.Reader
}

// Mock Close from readerMock (ReadCloser)
func (readerMock) Close() error {
	return nil
}
