package gara

import "net/http"

type HandlerFunc func(ctx *Context)

type Engine struct {
	router *router
}


func NewEngine() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := newContext(writer, request)
	e.router.handle(ctx)
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.addRoute("GET", path, handler)
}

func (e *Engine) POST(path string, handler HandlerFunc) {
	e.router.addRoute("POST", path, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}




