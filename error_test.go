/**
 * Created by zc on 2020/7/11.
 */
package validator

import (
	"errors"
	"testing"
)

func TestErrorChains_Error(t *testing.T) {
	tests := []struct {
		name string
		e    ErrorChains
		want string
	}{
		{
			name: "test1",
			e: ErrorChains{
				Error{
					path: "name",
					err: errors.New("not required"),
				},
				Error{
					path: "age",
					err: errors.New("not required"),
				},
			},
			want: "name: not required\nage: not required\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		path string
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				path: "name",
				err:  errors.New("name is not required"),
			},
			want: "name is not required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				path: tt.fields.path,
				err:  tt.fields.err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}