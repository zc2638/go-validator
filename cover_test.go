/**
 * Created by zc on 2020/7/11.
 */
package validator

import (
	"encoding/json"
	"testing"
)

type User struct {
	Name    string
	Age     int
	Address string
}

func TestJSONCover(t *testing.T) {
	user := User{
		Name:    "张三",
		Age:     18,
		Address: "北京市",
	}
	data, err := json.Marshal(&user)
	if err != nil {
		t.Errorf("JSONCover() json format fail")
	}
	tests := []struct {
		name    string
		data    []byte
		s       interface{}
		wantErr bool
	}{
		{
			name:    "test1",
			data:    data,
			s:       User{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := JSONCover()(tt.data, &tt.s); (err != nil) != tt.wantErr {
				t.Errorf("JSONCover() = %+v, want %v", tt.s, tt.wantErr)
			}
		})
	}
}
