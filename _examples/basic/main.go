package main

import (
	"fmt"
	"github.com/zc2638/go-validator"
)

/**
 * Created by zc on 2019-08-21.
 */

type User struct {
	Name interface{}            `vdr:"required,msg=姓名,必填" json:"name"`
	Age  int                    `vdr:"required,max=22,msg=年龄不对" json:"age"`
	M    map[string]interface{} `vdr:"required"`
	Addr Addr
}

type Addr struct {
	Name string `vdr:"required"`
}

func main() {

	num := 10
	str := "Hello World!"
	mapValue := map[string]string{}
	slice := make([]string, 0)

	vdr := validator.NewVdr()
	vdr.MakeValue(num, "max=15", "msg=超过最大值")
	vdr.MakeValue(str, "required,msg=字符串为空")
	vdr.MakeValue("asd2", "reg=^[a-z]*$")
	fmt.Println(vdr.Check())

	vdr.MakeValue(slice, "required")
	vdr.MakeValue(mapValue, "required,msg=map为空")
	fmt.Println(vdr.Check())

	var user = User{
		Name: "张三",
		Age:  18,
		M: map[string]interface{}{
			"test": "Hello",
		},
		Addr: Addr{
			Name: "北京市",
		},
	}
	fmt.Println(vdr.MakeStruct(user).Check())
}
