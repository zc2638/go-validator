/**
 * Created by zc on 2020/6/25.
 */
package validator

import (
	"bytes"
	"strings"
)

func buildSlicePath(path string) string {
	paths := strings.Split(path, ".")
	var buffer bytes.Buffer
	for _, p := range paths {
		index := strings.Index(p, SignSlice)
		if index == 0 {
			p = SignSlice
		}
		buffer.WriteString(p)
		buffer.WriteString(".")
	}
	bs := buffer.Bytes()
	pathBs := bs[:len(bs)-1]
	return string(pathBs)
}

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
