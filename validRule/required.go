package validRule

import (
	"github.com/zc2638/go-validator/typ"
	"strings"
)

/**
 * Created by zc on 2019-08-12.
 */

// 验证字符串是否为空
type Required struct {}

func (v *Required) Name() string                      { return "required" }
func (v *Required) SetCondition(...interface{}) error { return nil }
func (v *Required) Valid(val interface{}) error       {
	tc := typ.NewTypeC(val, typ.String)
	res, err := tc.Convert()
	if err != nil {
		return err
	}
	if strings.TrimSpace(res.(string)) == "" {
		return typ.NotRequired
	}
	return nil
}