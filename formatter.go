/**
 * Created by zc on 2020/6/12.
 */
package validator

import (
	"encoding/json"
	"reflect"
	"strconv"
)

func JSONFormatter() Formatter {
	return func(data []byte) (Current, error) {
		var res interface{}
		if err := json.Unmarshal(data, &res); err != nil {
			return nil, err
		}
		rv := reflect.ValueOf(res)
		current := make(Current)
		parse("", rv, current)
		return current, nil
	}
}

func parse(path string, rv reflect.Value, res map[string]interface{}) {
	if path != "" {
		res[path] = rv.Interface()
	}
CHECK:
	switch rv.Kind() {
	case reflect.Slice:
		if rv.Len() == 0 {
			return
		}
		for i := 0; i < rv.Len(); i++ {
			iv := strconv.Itoa(i)
			iPath := buildPath(path, SignSlice+ iv)
			parse(iPath, rv.Index(i), res)
		}
	case reflect.Map:
		if rv.Len() == 0 {
			return
		}
		keys := rv.MapKeys()
		for _, key := range keys {
			keyPath := buildPath(path, key.String())
			parse(keyPath, rv.MapIndex(key), res)
		}
	case reflect.Interface:
		rv = rv.Elem()
		goto CHECK
	}
}