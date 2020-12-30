package hiWeb

import (
	"log"
	"path"
)

type RouterGroup struct {
	prefix      string
	middlewares MiddlewareList
	engine      *Engine
}

func newRouterGroup(prefix string, engine *Engine) *RouterGroup {
	return &RouterGroup{
		prefix: prefix,
		engine: engine,
	}
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	return newRouterGroup(r.prefix+prefix, r.engine)
}

func (r *RouterGroup) Use(handlers ...HandlerFunc) {
	for _, v := range handlers {
		tmp := &Middleware{0, v}
		r.middlewares = append(r.middlewares, tmp)
	}
}

func (r *RouterGroup) UseWeight(middlewares ...Middleware) {
	for _, v := range middlewares {
		tmp := &Middleware{v.weight, v.handlerFunc}
		r.middlewares = append(r.middlewares, tmp)
	}
}

func (r *RouterGroup) addRoute(method string, str string, handler HandlerFunc) {
	pattern := r.prefix + str
	log.Printf("Route %4s - %s", method, pattern)
	r.engine.router.addRouter(method, pattern, handler)
}

func (r *RouterGroup) GET(pattern string, handler HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *RouterGroup) POST(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

func (r *RouterGroup) Static(relativePath, localPath string) {
	completePath := path.Join(r.prefix, relativePath)
	handler := staticFileHandler(completePath, localPath)
	urlPattern := path.Join(relativePath, "/*filepath")
	r.GET(urlPattern, handler)
}
