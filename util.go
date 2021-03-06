/**
 * Created by zc on 2020/7/5.
 */
package validator

import (
	"bytes"
	"reflect"
	"strings"
)

// value returns the reflect Value of val
func Value(val interface{}) reflect.Value {
	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv
}

func camelToUnderline(s string) string {
	num := len(s)
	data := make([]byte, 0, num*2)
	j := false
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func buildPath(paths ...string) string {
	var buffer bytes.Buffer
	for _, path := range paths {
		if path != "" {
			buffer.WriteString(path)
			buffer.WriteString(".")
		}
	}
	bs := buffer.Bytes()
	var pathBs string
	if len(bs) > 0 {
		pathBs = string(bs[:len(bs)-1])
	}
	return pathBs
}

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
