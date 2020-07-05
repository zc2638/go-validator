/**
 * Created by zc on 2020/6/26.
 */
package validator

import "reflect"

// typ returns the reflect Type of value
func Type(val interface{}) reflect.Type {
	rt := reflect.TypeOf(val)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt
}

// value returns the reflect Value of val
func Value(val interface{}) reflect.Value {
	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv
}
