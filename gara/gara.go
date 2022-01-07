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
	c := newContext(writer, request)
	e.router.handle(c)
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.addRoute("Get", path, handler)
}

func (e *Engine) POST(path string, handler HandlerFunc) {
	e.router.addRoute("Post", path, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}




