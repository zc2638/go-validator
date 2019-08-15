package validator

import (
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
