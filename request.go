package validator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

/**
 * Created by zc on 2019-08-14.
 */

type RequestContext struct {
	r *http.Request
}

func NewRequestContext(r *http.Request) *RequestContext {
	return &RequestContext{r: r}
}

func (c *RequestContext) Deadline() (deadline time.Time, ok bool) { return }
func (c *RequestContext) Done() <-chan struct{}                   { return nil }
func (c *RequestContext) Err() error                              { return nil }
func (c *RequestContext) Value(key interface{}) interface{} {
	req := c.r
	switch req.Method {
	case http.MethodGet:
		if values, ok := req.URL.Query()[key.(string)]; ok && len(values) > 0 {
			return values[0]
		}
	case http.MethodPost:
		return req.PostFormValue(key.(string))
	}
	return nil
}

// Request Body parser

/*
RequestBodyContext implement the context interface, and provide some parse functions.

this main to parse the request body, and for body stream can be open twice. it will not
break the stream of request stream, in other words, it copy the request stream for parse.

For use:
vdr := validator.SetContext(NewRequestBodyContext( { request } ))
res := vdr.CheckStruct( { your struct } )
*/
type RequestBodyContext struct {
	r         *http.Request
	bodyBytes [] byte
	readFunc  func(value []byte, targetKey interface{}) interface{}
}

/*
NewRequestBodyContext will create new RequestBodyContext, it only set http.Request.
*/
func NewRequestBodyContext(r *http.Request) *RequestBodyContext {
	return &RequestBodyContext{r: r}
}

func (c *RequestBodyContext) Deadline() (deadline time.Time, ok bool) { return }
func (c *RequestBodyContext) Done() <-chan struct{}                   { return nil }
func (c *RequestBodyContext) Err() error                              { return nil }
func (c *RequestBodyContext) Value(key interface{}) interface{} {
	// check parser func, if nil use default json parser func
	if c.readFunc == nil {
		c.readFunc = c.DefaultJsonParserFunc
	}
	if c.bodyBytes == nil || len(c.bodyBytes) == 0 {
		// get request body
		c.bodyBytes, _ = ioutil.ReadAll(c.r.Body)
		// close BufferReader to rewrite request body
		c.r.Body.Close()
		// rewrite request body
		c.r.Body = ioutil.NopCloser(bytes.NewBuffer(c.bodyBytes))
	}
	return c.readFunc(c.bodyBytes, key)
}

/*
SetReadFunc will set body parser func.

For different body, it may has different parse actions. So open the door of body parse func.
The parser func define:
	func FuncName(value []byte, targetKey interface{}) interface{} {
		func body...
	}

If doesn't invoke this func, RequestBodyContext will use default func :
	func (c RequestBodyContext) DefaultJsonParserFunc(value [] byte, targetKey interface{}) interface{}

And one kind of body only can use one kind of parser func. It means that you only write all operation of body's read and
parser. Else some operation will has error.

But some time, for different key may use different parse actions. It can write like this:
	func Func1(value []byte, targetKey interface{}) interface{} {
		func body...
	}
	func Func2(value []byte, targetKey interface{}) interface{} {
		func body...
	}
	func FuncAll(value []byte, targetKey interface{}) interface{} {
		if targetKey == "1"{
			return Func1(value, targetKey)
		}
		if targetKey == "2"{
			return Func2(value, targetKey)
		}

		func body...
	}
	======= set read func
	vdr := validator.SetContext(NewRequestBodyContext( { request } ).SetReadFunc(FuncAll))
	res := vdr.CheckStruct( { your struct } )
*/
func (c *RequestBodyContext) SetReadFunc(readFunc func(value []byte, targetKey interface{}) interface{}) *RequestBodyContext {
	c.readFunc = readFunc
	return c
}

/*
DefaultJsonParserFunc will set parser json string. It is the RequestBodyContext default read func.
Note: It only can parse only one floor json object. Such as { "a":"a"}, but can't parse : {"a":{"a":"a"}}.

If happen some errors, it will return the empty string( is "", not nil).
*/
func (c RequestBodyContext) DefaultJsonParserFunc(value [] byte, targetKey interface{}) interface{} {
	var data map[string]interface{}

	// convert json string to map[string]interface
	err := json.Unmarshal(value, &data)

	// has error will return ""
	if err != nil {
		return ""
	}

	return data[targetKey.(string)]
}
