package hiWeb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//保存原始的请求和响应
	writer http.ResponseWriter
	req    *http.Request

	//保存关键信息,用于快速寻找
	path   string //req.URL.Path
	method string //req.Method

	//url参数匹配
	params map[string]string

	//响应状态
	statusCode int

	//中间件 + 注册处理函数
	handlers      []HandlerFunc
	indexHandlers int

	//engine指针
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		writer:        w,
		req:           req,
		path:          req.URL.Path,
		method:        req.Method,
		params:        make(map[string]string),
		statusCode:    0,
		handlers:      make([]HandlerFunc, 0),
		indexHandlers: -1,
	}
}

func (c *Context) Next() {
	c.indexHandlers++
	s := len(c.handlers)
	for ; c.indexHandlers < s; c.indexHandlers++ {
		c.handlers[c.indexHandlers](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.indexHandlers = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value, _ := c.params[key]
	return value
}

func (c *Context) FormValue(key string) string {
	return c.req.FormValue(key)
}
func (c *Context) Query(key string) string {
	return c.req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.statusCode = code
	c.writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err.Error())
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExcuteTemplate(c.writer, name, data); err != nil {
		panic(err.Error())
	}
}
