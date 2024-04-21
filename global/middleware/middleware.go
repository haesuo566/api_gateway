package middleware

import (
	"net/http"
	"strings"
)

type middleware struct {
}

// Request Response Context
type Ctx struct {
	W http.ResponseWriter
	R *http.Request
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

type HandlerFunc func(*Ctx)

func New() *middleware {
	return &middleware{}
}

func (m *middleware) Get(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, GET, handler)
}

func (m *middleware) Post(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, POST, handler)
}

func (m *middleware) Put(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, PUT, handler)
}

func (m *middleware) Patch(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, PATCH, handler)
}

func (m *middleware) Delete(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, DELETE, handler)
}

func (m *middleware) Run(port string) error {
	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}

func setHandler(endPoint, methodType string, handler HandlerFunc) {
	http.HandleFunc(endPoint, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "" || !strings.EqualFold(methodType, r.Method) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		handler(&Ctx{w, r})
	})
}
