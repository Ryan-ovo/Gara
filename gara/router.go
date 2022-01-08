package gara

import (
	"net/http"
)

type router struct {
	// 根节点，每一种请求方式都有一棵路由树
	roots map[string]*trie
	// 处理函数
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*trie),
		handlers: make(map[string]HandlerFunc),
	}
}

// handle 处理路由
func (r *router) handle(ctx *Context) {
	node, param := r.getRoute(ctx.Method, ctx.Path)
	if node != nil {
		ctx.Params = param
		key := ctx.Method + "-" + node.path
		if handler, ok := r.handlers[key]; ok {
			handler(ctx)
		}
	} else {
		ctx.String(http.StatusNotFound, "404 not found: %s\n", ctx.Path)
	}
}

// addRoute 往路由树上插入路由
func (r *router) addRoute(method, path string, handler HandlerFunc) {
	// 获取这类请求方法（如get/post）的根节点，如果没有则创建
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &trie{son: make(map[string]*trie)}
	}
	root := r.roots[method]
	root.insert(path)
	// 存储每个路由的执行方法
	key := method + "-" + path
	r.handlers[key] = handler
}

// getRoute 从路由树上获取路由
func (r *router) getRoute(method, path string) (*trie, map[string]string) {
	var root *trie
	var ok bool
	if root, ok = r.roots[method]; !ok {
		return nil, nil
	}
	return root.search(path)
}


