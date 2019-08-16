package typ

import (
	"reflect"
)

/**
 * Created by zc on 2019-08-15.
 */

type typeS struct {
	t reflect.Type
	v reflect.Value
}

func NewTypeS(t reflect.Type, v reflect.Value) *typeS {
	return &typeS{t: t, v: v}
}

func (t *typeS) GetFiledSet() (fieldSet []reflect.StructField) {
	for i := 0; i < t.t.NumField(); i++ {
		fieldSet = append(fieldSet, t.t.Field(i))
	}
	return
}

func (t *typeS) SetValue(field reflect.StructField, value interface{}) error {
	v := t.v.FieldByName(field.Name)
	if v.CanSet() && value != nil {
		tc := NewTypeC(value, field.Type.Kind())
		res, err := tc.Convert()
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(res))
	}
	return nil
}
