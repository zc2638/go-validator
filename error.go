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
	e    error
}

func (e *Error) Error() string {
	return e.e.Error()
}

type ErrorChains []Error

func (e ErrorChains) Error() string {
	var buffer bytes.Buffer
	for _, err := range e {
		path := strings.ReplaceAll(err.path, SignSlice, "")
		buffer.WriteString(path)
		buffer.WriteString(": ")
		buffer.WriteString(err.e.Error())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
