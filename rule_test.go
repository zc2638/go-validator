/**
 * Created by zc on 2020/7/11.
 */
package validator

import (
	"reflect"
	"testing"
)

func TestRuleRequired(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		wantErr bool
	}{
		{
			name:    "test1",
			val:     nil,
			wantErr: true,
		},
		{
			name:    "test2",
			val:     "",
			wantErr: false,
		},
		{
			name:    "test3",
			val:     0,
			wantErr: false,
		},
		{
			name:    "test4",
			val:     "张三",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RuleRequired()(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("RuleRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRuleRequiredWithMessage(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		message string
		wantErr bool
	}{
		{
			name:    "test1",
			val:     nil,
			message: "not required",
			wantErr: true,
		},
		{
			name:    "test2",
			val:     "",
			message: "name is must",
			wantErr: false,
		},
		{
			name:    "test3",
			val:     0,
			message: "age is not required",
			wantErr: false,
		},
		{
			name:    "test4",
			val:     "张三",
			message: "name is not required",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RuleRequiredWithMessage(tt.message)(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("RuleRequiredWithMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRuleType(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		kind    reflect.Kind
		wantErr bool
	}{
		{
			name:    "test1",
			val:     float64(0),
			kind:    reflect.String,
			wantErr: true,
		},
		{
			name:    "test2",
			val:     "1",
			kind:    reflect.String,
			wantErr: false,
		},
		{
			name:    "test3",
			val:     float64(1),
			kind:    reflect.Int,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RuleType(tt.kind)(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("RuleType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
