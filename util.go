/**
 * Created by zc on 2020/7/5.
 */
package validator

import (
	"bytes"
	"strings"
)

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