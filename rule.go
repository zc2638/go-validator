package validator

import (
	"errors"
	"github.com/zc2638/go-validator/typ"
	"reflect"
	"strings"
)

/**
 * Created by zc on 2019-08-13.
 */

// 验证字符串是否为空
type RuleRequired struct{}

func (*RuleRequired) Name() string                      { return "required" }
func (*RuleRequired) SetCondition(...interface{}) error { return nil }
func (*RuleRequired) Fire(e *Engine) error {
	value := reflect.ValueOf(e.Part.Value)
	var errRequired = errors.New(e.Part.Key + "(" + value.Kind().String() + ") not required")

	if e.Part.Value == nil {
		return errRequired
	}
	switch value.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		if value.IsNil() {
			return errRequired
		}
	default:
		tc := typ.NewTypeC(e.Part.Value, reflect.String)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		if strings.TrimSpace(res.(string)) == "" {
			return errRequired
		}
	}
	return nil
}

// 验证类型
type RuleTypes struct {
	t string
}

func (r *RuleTypes) Name() string { return "type" }
func (r *RuleTypes) SetCondition(cs ...interface{}) error {
	t := typ.String
	if len(cs) > 0 {
		if reflect.TypeOf(cs[0]).Kind() != reflect.String {
			return typ.TypeNotString
		}
		t = cs[0].(string)
	}
	r.t = t
	return nil
}

func (r *RuleTypes) Fire(e *Engine) error {
	tc := typ.NewTypeC(e.Part.Value, typ.ChangeTypeToKind(r.t))
	_, err := tc.Convert()
	if err != nil {
		return err
	}
	return nil
}

// 验证最大值
type RuleMax struct {
	n float64
}

func (r *RuleMax) Name() string { return "max" }
func (r *RuleMax) SetCondition(cs ...interface{}) error {
	if len(cs) > 0 {
		tc := typ.NewTypeC(cs[0], reflect.Float64)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		r.n = res.(float64)
	}
	return nil
}

func (r *RuleMax) Fire(e *Engine) error {
	tc := typ.NewTypeC(e.Part.Value, reflect.Float64)
	res, err := tc.Convert()
	if err != nil {
		return err
	}

	if res.(float64) > r.n {
		return typ.NumberMaxLimit
	}
	return nil
}

// 验证最小值
type RuleMin struct {
	n float64
}

func (r *RuleMin) Name() string { return "min" }
func (r *RuleMin) SetCondition(cs ...interface{}) error {
	if len(cs) > 0 {
		tc := typ.NewTypeC(cs[0], reflect.Float64)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		r.n = res.(float64)
	}
	return nil
}

func (r *RuleMin) Fire(e *Engine) error {
	tc := typ.NewTypeC(e.Part.Value, reflect.Float64)
	res, err := tc.Convert()
	if err != nil {
		return err
	}

	if res.(float64) < r.n {
		return typ.NumberMinLimit
	}
	return nil
}

// 验证长度
type RuleLen struct {
	len int
}

func (r *RuleLen) Name() string { return "len" }
func (r *RuleLen) SetCondition(cs ...interface{}) error {
	if len(cs) > 0 {
		tc := typ.NewTypeC(cs[0], reflect.Int)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		r.len = res.(int)
	}
	return nil
}

func (r *RuleLen) Fire(e *Engine) error {
	// TODO
	return nil
}
