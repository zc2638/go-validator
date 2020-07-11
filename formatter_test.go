/**
 * Created by zc on 2020/7/11.
 */
package validator

import (
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	test1 := map[string]interface{}{
		"name":    "张三",
		"age":     18,
		"address": "北京市",
	}
	test2 := []map[string]interface{}{test1, test1}
	type args struct {
		path string
		rv   reflect.Value
		res  map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "test1",
			args: args{
				path: "",
				rv:   reflect.ValueOf(test1),
				res:  make(map[string]interface{}),
			},
			want: map[string]interface{}{
				"name":    "张三",
				"age":     18,
				"address": "北京市",
			},
		},
		{
			name: "test2",
			args: args{
				path: "",
				rv:   reflect.ValueOf(test2),
				res:  make(map[string]interface{}),
			},
			want: map[string]interface{}{
				SignSlice + "0":         test1,
				SignSlice + "0.name":    "张三",
				SignSlice + "0.age":     18,
				SignSlice + "0.address": "北京市",
				SignSlice + "1":         test1,
				SignSlice + "1.name":    "张三",
				SignSlice + "1.age":     18,
				SignSlice + "1.address": "北京市",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parse(tt.args.path, tt.args.rv, tt.args.res)
			if !reflect.DeepEqual(tt.args.res, tt.want) {
				t.Errorf("parse() = %v, want %v", tt.args.res, tt.want)
			}
		})
	}
}
