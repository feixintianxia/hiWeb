package hiWeb

import "strings"

type router struct {
	trieRoots      map[string]*trieNode   //key = ["GET","POST"]
	registHandlers map[string]HandlerFunc //注册函数
}

func newRouter() *router {
	return &router{
		trieRoots:      make(map[string]*trieNode),
		registHandlers: make(map[string]HandlerFunc),
	}
}

func getRegistHandlersKey(method, pattern string) string {
	return method + "-" + pattern
}

// /a/b/c => ["a", "b", "c"]
func parsePattern(pattern string) []string {
	arr := strings.Split(pattern, "/")

	parts := make([]string, 0)

	for _, item := range arr {
		if item != "" {
			parts = append(parts, item)
			//仅仅允许一个 *
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRouter(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.trieRoots[method]
	if !ok {
		r.trieRoots[method] = newRootNode()
	}

	r.trieRoots[method].insert(pattern, parts)
	r.registHandlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*trieNode, map[string]string) {
	orginParts := parsePattern(path)
	params := make(map[string]string)
	rootNode, ok := r.trieRoots[method]

	if !ok {
		return nil, nil
	}

	node := rootNode.search(orginParts)
	if node == nil {
		return nil, nil
	}

	patternParts := parsePattern(node.pattern)
	for index, part := range patternParts {
		if part[0] == ':' {
			params[part[1:]] = orginParts[index]
		}

		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(orginParts[index:], "/")
			break
		}
	}
	return node, params
}

func (r *router) getRoutes(method string) []*trieNode {
	rootNode, ok := r.trieRoots[method]
	if !ok {
		return nil
	}

	nodes := make([]*trieNode, 0)
	rootNode.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.method, c.path)

	if n == nil {
		c.handlers = append(c.handlers, NotFoundResponse)
	} else {
		c.params = params

		key := getRegistHandlersKey(c.method, n.pattern)
		c.handlers = append(c.handlers, r.registHandlers[key])
	}

	c.Next()
}
