/**
 * Created by zc on 2020/6/25.
 */
package validator

func MakeSliceValue(vfs ...ValidateFunc) HandlerFunc {
	return func(v Validation) {
		v.MakeField(SignSlice, vfs...)
	}
}

func MakeSliceHandler(handler Handler, vfs ...ValidateFunc) HandlerFunc {
	return func(v Validation) {
		v.Make(handler, vfs...)
	}
}

func MakeSliceSub(s interface{}, f HandlerFunc, vfs ...ValidateFunc) HandlerFunc {
	return func(v Validation) {
		v.MakeSlice(s, f, vfs...)
	}
}
