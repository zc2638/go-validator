package main

import (
	"fmt"
	"github.com/zc2638/go-validator"
	"net/http"
)

/**
 * Created by zc on 2019-08-21.
 */

type User struct {
	Name string `vdr:"required,msg=姓名,必填" json:"name"`
	Age  int    `vdr:"min=10,max=22,msg=年龄不对" json:"age"`
}

func main() {

	http.HandleFunc("/user/struct", func(writer http.ResponseWriter, request *http.Request) {
		var user User
		vdr := validator.SetContext(validator.NewRequestContext(request))
		fmt.Println(vdr.CheckStruct(&user))

		_, _ = writer.Write([]byte("Hello World!"))
	})

	http.HandleFunc("/user/slice", func(writer http.ResponseWriter, request *http.Request) {
		var vs = [][]string{
			{"name", `required,msg=姓名,必填`},
			{"age", `min=10,max=22`},
		}
		fmt.Println(validator.SetContext(validator.NewRequestContext(request)).CheckSlice(vs...))

		_, _ = writer.Write([]byte("Hello World!"))
	})

	http.HandleFunc("/user/map", func(writer http.ResponseWriter, request *http.Request) {
		var vm = map[string]string{
			"name": `required,msg=姓名,必填`,
			"age":  `min=10,max=22`,
		}
		fmt.Println(validator.SetContext(validator.NewRequestContext(request)).CheckMap(vm))

		_, _ = writer.Write([]byte("Hello World!"))
	})

	_ = http.ListenAndServe(":8080", nil)
}
