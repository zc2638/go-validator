/**
 * Created by zc on 2020/6/25.
 */
package validator

import "bytes"

type Error struct {
	path string
	e    error
}

type ErrorChains []Error

func (e ErrorChains) Error() string {
	var buffer bytes.Buffer
	for _, err := range e {
		buffer.WriteString(err.path)
		buffer.WriteString(": ")
		buffer.WriteString(err.e.Error())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
