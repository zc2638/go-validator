package validRule

import (
	"reflect"
	"validator/typ"
)

/**
 * Created by zc on 2019-08-12.
 */

// 验证类型
type Types struct {
	t string
}

func (v *Types) Name() string { return "type" }
func (v *Types) SetCondition(cs ...interface{}) error {
	t := typ.String
	if len(cs) > 0 {
		if reflect.TypeOf(cs[0]).Kind() != reflect.String {
			return typ.TypeNotString
		}
		t = cs[0].(string)
	}
	v.t = t
	return nil
}

func (v *Types) Valid(val interface{}) error {
	tc := typ.NewTypeC(val, v.t)
	_, err := tc.Convert()
	if err != nil {
		return err
	}
	return nil
}
