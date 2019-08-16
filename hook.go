package validator

import (
	"errors"
	"github.com/zc2638/go-validator/typ"
	"reflect"
)

/**
 * Created by zc on 2019-08-13.
 */

type HookMsg struct {
	s string
}

func (h *HookMsg) Name() string {
	return "msg"
}

func (h *HookMsg) SetCondition(cs ...interface{}) error {
	if len(cs) > 0 {
		tc := typ.NewTypeC(cs[0], reflect.String)
		if !tc.CheckKind() {
			return typ.TypeNotString
		}
		h.s = cs[0].(string)
	}
	return nil
}

func (h *HookMsg) Fire(e *Engine) error {
	if e.Err != nil {
		e.Err = errors.New(h.s)
	}
	return nil
}