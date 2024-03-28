package ancho

import (
	"net/http"
	"strings"
)

type ancho struct {
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

func New() *ancho {
	return &ancho{}
}

func (a *ancho) Get(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, GET, handler)
}

func (a *ancho) Post(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, POST, handler)
}

func (a *ancho) Put(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, PUT, handler)
}

func (a *ancho) Patch(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, PATCH, handler)
}

func (a *ancho) Delete(endPoint string, handler HandlerFunc) {
	setHandler(endPoint, DELETE, handler)
}

func (a *ancho) Run(port string) error {
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
