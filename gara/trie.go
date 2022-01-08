package gara

import "strings"

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

func (t *trie) insert(path string) {
	parts := parsePath(path)
	p := t
	for _, part := range parts {
		// 如果子节点不存在直接创建
		// 子节点的part记录了当前节点的单词，子节点的path记录了从根节点走到当前节点的全路径
		// 标记该节点是否是通配符节点
		if p.son[part] == nil {
			p.son[part] = &trie{
				part:   part,
				son:    make(map[string]*trie),
				isWild: part[0] == '*' || part[0] == ':',
			}
		}
		p = p.son[part]
	}
	p.path = path
}

func (t *trie) search(path string) (*trie, map[string]string) {
	parts := parsePath(path)
	param := make(map[string]string)
	p := t
	for i, part := range parts {
		var temp string
		for _, node := range p.son {
			// 遍历每个子节点，匹配整个单词或者通配符
			if node.part == part || node.isWild {
				// 如果匹配到的是通配符，则把通配符后的变量赋值成匹配到的真实变量
				if node.part[0] == '*' {
					param[node.part[1:]] = strings.Join(parts[i:], "/")
				} else if node.part[0] == ':' {
					param[node.part[1:]] = part
				}
				temp = node.part
			}
		}
		// 如果匹配到*说明后面不需要再匹配了，直接返回即可
		if temp[0] == '*' {
			return p.son[temp], param
		}
		// 走到匹配上的节点进行下一轮的匹配
		p = p.son[temp]
	}
	return p, param
}

// 解析路径，只允许存在一个*
// /hello/a/b解析为[hello, a, b]
// /hello/:a/b解析为[hello, :a, b]
// /hello/*a/b解析为[hello, *a]
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
