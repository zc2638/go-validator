package validRule

import (
	"github.com/zc2638/go-validator/typ"
)

/**
 * Created by zc on 2019-08-12.
 */

type Max struct {
	n float64
}

func (v *Max) Name() string { return "max" }

func (v *Max) SetCondition(cs ...interface{}) error {
	if len(cs) > 0 {
		tc := typ.NewTypeC(cs[0], typ.Float64)
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		v.n = res.(float64)
	}
	return nil
}

func (v *Max) Valid(val interface{}) error {
	tc := typ.NewTypeC(val, typ.Float64)
	res, err := tc.Convert()
	if err != nil {
		return err
	}

	if res.(float64) > v.n {
		return typ.NumberMaxLimit
	}
	return nil
}
