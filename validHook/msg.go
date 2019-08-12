package validHook

import (
	"errors"
	"github.com/zc2638/go-validator"
	"github.com/zc2638/go-validator/typ"
)

/**
 * Created by zc on 2019-08-12.
 */

type Msg struct {
	s string
}

func (h *Msg) Name() string {
	return "msg"
}

func (h *Msg) SetCondition(cs ...interface{}) error {
	if len(cs) > 0 {
		tc := typ.NewTypeC(cs[0], typ.String)
		if !tc.CheckKind() {
			return typ.TypeNotString
		}
		h.s = cs[0].(string)
	}
	return nil
}

func (h *Msg) Fire(e *validator.VdrEngine) error {
	if e.Err != nil {
		e.Err = errors.New(h.s)
	}
	return nil
}