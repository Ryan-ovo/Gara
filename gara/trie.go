package gara

type trie struct {
	// 全路由
	path   string
	// 路由的每个部分
	part   string
	// 子节点
	son    map[string]*trie
	// 是否为通配符节点
	isWild bool
}

