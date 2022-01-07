package gara

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// 请求和响应
	Writer http.ResponseWriter
	Req    *http.Request
	// 路由信息
	Path   string
	Method string
	// 响应码
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
	}
}

// PostForm acquire specific field from form
// 获取表达你对应字段的值
func (c *Context) PostForm(key string) string {
	// acquire the first field from Request.Form
	return c.Req.FormValue(key)
}

// Query 获取参数对应字段的值
// acquire first value use r.URL.Query().Get("ParamName")
// acquire all value use r.URL.Query()["ParamName"]
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status store status code
func (c *Context) Status(code int) {
	c.StatusCode = code
	// WriteHeader sends an HTTP response header with the provided status code.
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	// A Header represents the key-value pairs in an HTTP header.
	c.Writer.Header().Set(key, value)
}

// String 封装String类型的响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 封装Json类型的响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 封装Data类型的响应
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 封装Html类型的响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
