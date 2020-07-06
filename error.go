/**
 * Created by zc on 2020/6/25.
 */
package validator

import (
	"bytes"
	"strings"
)

type Error struct {
	path string
	err  error
}

func (e *Error) Error() string {
	return e.err.Error()
}

type ErrorChains []Error

func (e ErrorChains) Error() string {
	var buffer bytes.Buffer
	for _, err := range e {
		path := strings.ReplaceAll(err.path, SignSlice, "")
		buffer.WriteString(path)
		buffer.WriteString(": ")
		buffer.WriteString(err.err.Error())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
