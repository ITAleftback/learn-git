package gin

import (
	"fmt"
	"net/http"
	"strings"
)

type Context struct {
	req *http.Request
	w    http.ResponseWriter
	queryParam map[string]string
	formParam map[string]string
	Writer    ResponseWriter
}
//来个 定义gin.H
type H map[string]interface{}


func (c *Context)String(s string){
	_,_=c.Writer.Write([]byte(s))
}

//写的辣鸡json     也不知道对不对===============================================
//这里很多都抄了gin框架
type JSON struct {
	Data interface{}
}

func (J JSON) Render(http.ResponseWriter) error {
	panic("implement me")
}

func (J JSON) WriteContentType(w http.ResponseWriter) {
	panic("implement me")
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Render(code, JSON{Data: obj})
}
func (c *Context) Render(code int, r Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeaderNow()
		return
	}
	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}
type Render interface {
	// Render writes data with custom ContentType.
	Render(http.ResponseWriter) error
	// WriteContentType writes custom ContentType.
	WriteContentType(w http.ResponseWriter)
}
func (c *Context) Status(code int) {
	c.w.WriteHeader(code)
}
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// Returns the HTTP response status code of the current request.
	Status() int

	// Returns the number of bytes already written into the response http body.
	// See Written()
	Size() int

	// Writes the string into the response body.
	WriteString(string) (int, error)

	// Returns true if the response body was already written.
	Written() bool

	// Forces to write the http header (status code + headers).
	WriteHeaderNow()

	// get the http.Pusher for server push
	Pusher() http.Pusher
}
//=========================================================================================





func NewContext(rw http.ResponseWriter,r *http.Request)(ctx Context){
	ctx=Context{
		req:         r,
		w:          rw,
		formParam:  make(map[string]string),
	}
	ctx.queryParam=parseQuery(r.RequestURI)
	return
}

func (c *Context)Query(key string)string{
	v:=c.queryParam[key]
	return  v
}

func parseQuery(uri string)(res map[string]string){
	res=make(map[string]string)
	uris :=strings.Split(uri,"?")
	if len(uris)==1{
		return
	}
	param:=uris[len(uris)-1]
	pair :=strings.Split(param,"&")
	for _,kv:=range pair{
		kvPair :=strings.Split(kv,"=")
		if len(kvPair)!=2{
			fmt.Println(kvPair)
			panic("request error")
		}
		res[kvPair[0]]=kvPair[1]
	}
	return
}
