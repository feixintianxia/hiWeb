package hiWeb

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// 时间统计
func TimeStatic() HandlerFunc {
	return func(c *Context) {
		begin := time.Now()
		c.Next()
		end := time.Now()
		delta := end.Sub(begin)
		log.Printf("%s-%s, delta: %d ms, req: %s, statuscode: %d",
			begin.Format(TIMEFORMAT),
			end.Format(TIMEFORMAT),
			delta,
			c.statusCode,
			c.req.RequestURI,
		)
	}
}

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
