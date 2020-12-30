package hiWeb

import "net/http"

const (
	TIMEFORMAT = "2006-01-02 15:04:05"
)

type HandlerFunc func(*Context)

func NotFoundResponse(c *Context) {
	c.String(http.StatusNotFound, "NOT FOUND 404: %s\n", c.path)
}

//urlPath 是 请求过来的url
//localPath 是注册时候的本地路径
func staticFileHandler(urlPath string, localPath string) HandlerFunc {
	return func(c *Context) {
		fs := http.FileSystem(http.Dir(localPath))
		filePath := c.Param("filepath")
		if _, err := fs.Open(filePath); err != nil {
			NotFoundResponse(c)
			return
		}

		fileServer := http.StripPrefix(urlPath, http.FileServer(fs))
		fileServer.ServeHTTP(c.writer, c.req)
	}
}

//可定制权重的中间件
type Middleware struct {
	weight      int //越小越也靠前
	handlerFunc HandlerFunc
}

type MiddlewareList []*Middleware

func (m MiddlewareList) Len() int           { return len(m) }
func (m MiddlewareList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MiddlewareList) Less(i, j int) bool { return m[i].weight < m[j].weight }

func (m MiddlewareList) getHandlers() []HandlerFunc {
	var handlers []HandlerFunc
	for _, v := range m {
		handlers = append(handlers, v.handlerFunc)
	}
	return handlers
}
