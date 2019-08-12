package validRule

import "github.com/zc2638/go-validator/typ"

/**
 * Created by zc on 2019-08-12.
 */

type Min struct {
	n float64
}

func (v *Min) Name() string { return "min" }

func (v *Min) SetCondition(cs ...interface{}) error {
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

func (v *Min) Valid(val interface{}) error {
	tc := typ.NewTypeC(val, typ.Float64)
	res, err := tc.Convert()
	if err != nil {
		return err
	}

	if res.(float64) < v.n {
		return typ.NumberMinLimit
	}
	return nil
}
