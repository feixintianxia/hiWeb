package hiWeb

import (
	"html/template"
	"net/http"
	"sort"
	"strings"
)

//web引擎
type Engine struct {
	*RouterGroup
	router        *router            //负责注册和查找路由
	groups        []*RouterGroup     //路由分组
	htmlTemplates *template.Template //html模板
	htmlFuncMap   template.FuncMap   //html模板自定义处理函数
}

func NewEngine() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = newRouterGroup("", engine)
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := NewEngine()
	engine.Use(TimeStatic(), Recovery())
	return engine
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.htmlFuncMap = funcMap
}

func (e *Engine) LoadHtmlTemplates(pattern string) {
	e.htmlTemplates = template.Must(
		template.New("hiWeb").
			Funcs(e.engine.htmlFuncMap).
			ParseGlob(pattern))
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	middlewares := make(MiddlewareList, 0)
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	sort.Sort(middlewares)

	c := newContext(w, req)
	c.handlers = middlewares.getHandlers()
	c.engine = e
	e.router.handle(c)
}
