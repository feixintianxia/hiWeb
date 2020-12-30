package hiWeb

import (
	"fmt"
	"strings"
)

type trieNode struct {
	pattern  string      // 注册的模式，不会是所有的节点都会保存 /a/b/c
	part     string      // 一个路径 /a
	children []*trieNode //子节点
	isWild   bool        //当前节点是否模糊节点 /:name 或 /*filepath
	height   int         //树高
}

func newRootNode() *trieNode {
	return &trieNode{
		height: 0,
	}
}

func newChildNode(part string, height int) *trieNode {
	return &trieNode{
		part:   part,
		isWild: part[0] == ':' || part[0] == '*',
		height: height,
	}
}

func (t *trieNode) String() string {
	return fmt.Sprintf(`trieNode{
		pattern: \"%s\",  
		part: \"%s\",
		isWild: %t,
		height: %d}`,
		t.pattern, t.part, t.isWild, t.height)
}

func (t *trieNode) matchOneChild(part string) *trieNode {
	var node *trieNode = nil
	for _, child := range t.children {
		if child.part == part || child.isWild {
			node = child
			break
		}
	}
	return node
}

func (t *trieNode) matchAllChild(part string) []*trieNode {
	nodes := make([]*trieNode, 0)
	for _, child := range t.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 如果pattern: /a/b/c, 则 parts: ["a", "b", "c"]
// 根节点 height = 0    part = ""
// 第一层节点 height = 1 part = "a"
func (t *trieNode) insert(pattern string, parts []string) {
	//先判断自身是否满足条件
	if t.height == len(parts) {
		t.pattern = pattern
		return
	}

	//取下一个路径
	part := parts[t.height]
	child := t.matchOneChild(part)
	if child == nil {
		child = newChildNode(part, t.height+1)
		t.children = append(t.children, child)
	}
	child.insert(pattern, parts)
}

func (t *trieNode) search(parts []string) *trieNode {
	if t.height+1 >= len(parts) || strings.HasPrefix(t.part, "*") {
		if t.pattern == "" {
			return nil
		}
		return t
	}

	//取下一个路径
	childHeight := t.height + 1
	part := parts[childHeight]
	childrens := t.matchAllChild(part)

	var node *trieNode = nil
	//DFS
	for _, child := range childrens {
		if node = child.search(parts); node != nil {
			break
		}
	}
	return node
}

func (t *trieNode) travel(nodes *([]*trieNode)) {
	if t.pattern != "" {
		*nodes = append(*nodes, t)
	}

	for _, child := range t.children {
		child.travel(nodes)
	}
}
