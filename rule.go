package validator

import (
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
func (*RuleRequired) Fire(e *VdrEngine) error {
	if e.Val == nil {
		return typ.NotRequired
	}
	tc := typ.NewTypeC(e.Val, typ.String)
	res, err := tc.Convert()
	if err != nil {
		return err
	}
	if strings.TrimSpace(res.(string)) == "" {
		return typ.NotRequired
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

func (r *RuleTypes) Fire(e *VdrEngine) error {
	tc := typ.NewTypeC(e.Val, r.t)
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
		tc := typ.NewTypeC(cs[0], typ.Float64)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		r.n = res.(float64)
	}
	return nil
}

func (r *RuleMax) Fire(e *VdrEngine) error {
	tc := typ.NewTypeC(e.Val, typ.Float64)
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
		tc := typ.NewTypeC(cs[0], typ.Float64)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		r.n = res.(float64)
	}
	return nil
}

func (r *RuleMin) Fire(e *VdrEngine) error {
	tc := typ.NewTypeC(e.Val, typ.Float64)
	res, err := tc.Convert()
	if err != nil {
		return err
	}

	if res.(float64) < r.n {
		return typ.NumberMinLimit
	}
	return nil
}