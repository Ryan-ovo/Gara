package gara

import (
	"fmt"
	"net/http"
	"strings"
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

func (r *router) handle(ctx *Context) {
	node, param := r.getRoute(ctx.Method, ctx.Path)
	for k, v := range r.handlers {
		fmt.Println(k, v)
	}
	if node != nil {
		ctx.Params = param
		key := ctx.Method + "-" + node.path
		if handler, ok := r.handlers[key]; ok {
			handler(ctx)
		}else {
			fmt.Println("hahaha")
		}
	} else {
		ctx.String(http.StatusNotFound, "404 not found: %s\n", ctx.Path)
	}
}

func (r *router) addRoute(method, path string, handler HandlerFunc) {
	parts := parsePath(path)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &trie{son: make(map[string]*trie)}
	}
	root := r.roots[method]
	key := method + "-" + path
	for _, part := range parts {
		if root.son[part] == nil {
			root.son[part] = &trie{
				part:   part,
				son:    make(map[string]*trie),
				isWild: part[0] == '*' || part[0] == ':',
			}
		}
		root = root.son[part]
	}
	root.path = path
	r.handlers[key] = handler
}

func (r *router) getRoute(method, path string) (*trie, map[string]string) {
	parts := parsePath(path)
	param := make(map[string]string)
	var root *trie
	var ok bool
	if root, ok = r.roots[method]; !ok {
		return nil, nil
	}
	for i, part := range parts {
		var temp string
		for _, node := range root.son {
			if node.part == part || node.isWild {
				if node.part[0] == '*' {
					param[node.part[1:]] = strings.Join(parts[i:], "/")
				} else if node.part[0] == ':' {
					param[node.part[1:]] = part
				}
				temp = node.part
			}
		}
		if temp[0] == '*' {
			return root.son[temp], param
		}
		root = root.son[temp]
	}
	return root, param
}

func (r *router) getRoute2(method, path string) (*trie, map[string]string) {
	parts := parsePath(path)
	param := make(map[string]string)
	var root *trie
	var ok bool
	if root, ok = r.roots[method]; !ok {
		return nil, nil
	}
	for i, part := range parts {
		if root.son[part] == nil {
			return nil, nil
		}
		node := root.son[part]
		if node.isWild {
			if node.part[0] == '*' {
				param[node.part[1:]] = strings.Join(parts[i:], "/")
			} else if node.part[0] == ':' {
				param[node.part[1:]] = part
			}
		}
		if node.part[0] == '*' {
			return root, param
		}
		root = root.son[part]
	}
	return root, param
}

func parsePath(path string) []string {
	res := make([]string, 0)
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part != "" {
			res = append(res, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return res
}
