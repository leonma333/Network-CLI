package network

import "net/http"

// HttpServer interface to define some http server functions
type HttpServer interface {
	ListenAndServe(addr string, handler http.Handler) error
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Get(url string) (resp *http.Response, err error)
}

type httpServer struct{}

func (h *httpServer) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

func (h *httpServer) Handle(pattern string, handler http.Handler) {
	http.Handle(pattern, handler)
}

func (h *httpServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}

func (h *httpServer) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
}
