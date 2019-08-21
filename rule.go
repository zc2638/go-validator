package validator

import (
	"errors"
	"github.com/zc2638/go-validator/typ"
	"reflect"
	"regexp"
	"strings"
)

/**
 * Created by zc on 2019-08-13.
 */

// 验证字符串是否为空
type RuleRequired struct{}

func (*RuleRequired) Name() string { return "required" }
func (*RuleRequired) Fire(e *Engine) error {
	value := e.Part.Value
	var errRequired = errors.New(e.Part.Key + "(" + value.Kind().String() + ") not required")
	switch value.Kind() {
	case reflect.Invalid:
		e.Err = errRequired
	case reflect.Map:
		if len(value.MapKeys()) == 0 {
			e.Err = errRequired
		}
	case reflect.Slice, reflect.Array:
		if value.Len() == 0 {
			e.Err = errRequired
		}
	default:
		tc := typ.NewTypeC(value.Interface(), reflect.String)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		if strings.TrimSpace(res.(string)) == "" {
			e.Err = errRequired
		}
	}
	return nil
}

// 验证最大值
type RuleMax struct{}

func (*RuleMax) Name() string { return "max" }
func (*RuleMax) Fire(e *Engine) error {

	required := new(RuleRequired)
	if err := required.Fire(e); err != nil {
		return err
	}
	if e.Err != nil {
		return nil
	}

	var max float64
	if e.Condition != "" {
		tc := typ.NewTypeC(e.Condition, reflect.Float64)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		max = res.(float64)
	}

	tc := typ.NewTypeC(e.Part.Value.Interface(), reflect.Float64)
	res, err := tc.Convert()
	if err != nil {
		return err
	}

	if res.(float64) > max {
		e.Err = errors.New(e.Part.Key + "(" + e.Part.Value.Kind().String() + ") is over max limit")
	}
	return nil
}

// 验证最小值
type RuleMin struct{}

func (*RuleMin) Name() string { return "min" }
func (*RuleMin) Fire(e *Engine) error {

	required := new(RuleRequired)
	if err := required.Fire(e); err != nil {
		return err
	}
	if e.Err != nil {
		return nil
	}

	var min float64
	if e.Condition != "" {
		tc := typ.NewTypeC(e.Condition, reflect.Float64)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		min = res.(float64)
	}

	tc := typ.NewTypeC(e.Part.Value.Interface(), reflect.Float64)
	res, err := tc.Convert()
	if err != nil {
		return err
	}

	if res.(float64) < min {
		e.Err = errors.New(e.Part.Key + "(" + e.Part.Value.Kind().String() + ") is below min limit")
	}
	return nil
}

// 验证长度
type RuleLen struct{}

func (*RuleLen) Name() string { return "len" }
func (*RuleLen) Fire(e *Engine) error {
	required := new(RuleRequired)
	if err := required.Fire(e); err != nil {
		return err
	}
	if e.Err != nil {
		return nil
	}

	var length int
	if e.Condition != "" {
		tc := typ.NewTypeC(e.Condition, reflect.Int)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		length = res.(int)
	}

	var num int
	switch e.Part.Value.Kind() {
	case reflect.Map:
		num = len(e.Part.Value.MapKeys())
	case reflect.Slice, reflect.Array:
		num = e.Part.Value.Len()
	default:
		tc := typ.NewTypeC(e.Part.Value.Interface(), reflect.String)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		num = strings.Count(res.(string), "") - 1
	}

	if num != length {
		e.Err = errors.New(e.Part.Key + "(" + e.Part.Value.Kind().String() + ") length not match")
	}
	return nil
}

// 正则
type RuleRegexp struct{}

func (*RuleRegexp) Name() string { return "reg" }
func (*RuleRegexp) Fire(e *Engine) error {
	required := new(RuleRequired)
	if err := required.Fire(e); err != nil {
		return err
	}
	if e.Err != nil {
		return nil
	}
	if e.Part.Value.Kind() != reflect.String {
		return errors.New(e.Part.Key + "(" + e.Part.Value.Kind().String() + ") is not string")
	}

	if e.Condition != "" {
		reg, err := regexp.Compile(e.Condition)
		if err != nil {
			return err
		}
		if !reg.MatchString(e.Part.Value.String()) {
			e.Err = errors.New(e.Part.Key + "(" + e.Part.Value.Kind().String() + ") not match")
		}
	}
	return nil
}
