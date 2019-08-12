package validRule

import (
	"github.com/zc2638/go-validator/typ"
	"regexp"
)

/**
 * Created by zc on 2019-08-12.
 */

type Regexp struct {
	rs []*regexp.Regexp
}

func (v *Regexp) Name() string { return "regexp" }

func (v *Regexp) SetCondition(cs ...interface{}) error {
	if len(cs) > 0 {
		if v.rs == nil {
			v.rs = make([]*regexp.Regexp, 0)
		}
		for _, c := range cs {
			reg, err := regexp.Compile(c.(string))
			if err != nil {
				return err
			}
			v.rs = append(v.rs, reg)
		}
	}
	return nil
}

func (v *Regexp) Valid(val interface{}) error {

	tc := typ.NewTypeC(val, typ.String)
	res, err := tc.Convert()
	if err != nil {
		return err
	}
	for _, r := range v.rs {
		if !r.MatchString(res.(string)) {
			return typ.RegexpNotMatch
		}
	}
	return nil
}
