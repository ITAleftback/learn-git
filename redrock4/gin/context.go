package gin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Context struct {
	req *http.Request
	w    http.ResponseWriter
	queryParam map[string]string
	formParam map[string]string
}
//来个 定义gin.H
type H map[string]interface{}

type Json struct {
	Data interface{}
}
func (c *Context)String(s string){
	_,_=c.w.Write([]byte(s))
}

//写的辣鸡json  就用了个数据序列化   也不知道对不对
func (c *Context)Json(message interface{})([]byte){
	jsonmessages,_:=json.Marshal(&Json{Data:message})
	return jsonmessages
}
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
