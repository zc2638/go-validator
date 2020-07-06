package main

import (
	"fmt"
	"github.com/zc2638/go-validator"
)

/**
 * Created by zc on 2019-08-21.
 */

type User struct {
	Name  interface{} `json:"name" vdr:"required,msg=姓名,必填"`
	Age   int         `json:"age" vdr:"required,max=22,msg=年龄不对"`
	M     map[string]interface{}
	Addr  Addr   `json:"addr"`
	Cates []Cate `json:"cates"`
}

func (u *User) Validate(validate validator.Validation) {
	validate.MakeValue(&u.Name, validator.RuleRequiredWithMessage("name必填"))
	validate.MakeValue(&u.Age)
	validate.MakeValue(&u.M)
	validate.Make(&u.Addr, validator.RuleRequiredWithMessage("addr为空"))
	validate.MakeSlice(&u.Cates, validator.MakeSliceHandler(&Cate{}), validator.RuleRequiredWithMessage("cate为空"))
}

func (u *User) Say() {
	fmt.Println("Hello World!")
}

type Addr struct {
	Name  string `vdr:"required"`
	Point Point
}

func (a *Addr) Validate(validate validator.Validation) {
	validate.MakeValue(&a.Name)
	validate.Make(&a.Point)
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p *Point) Validate(validate validator.Validation) {
	validate.MakeValue(&p.X, validator.RuleRequiredWithMessage("x坐标为空"))
	validate.MakeValue(&p.Y)
}

type Cate struct {
	Name string `json:"name"`
}

func (c *Cate) Validate(validate validator.Validation) {
	validate.MakeValue(&c.Name, validator.RuleRequiredWithMessage("cate name 为空"))
}

var str = `{
  "name": "张三",
  "age": 25,
  "m": {},
  "addr": {
  	"name": "上海",
	"point": {
		"x": 12
	}
  },
  "cates": [
  	{ "name": "123" }
  ]
}`

func main() {
	var user User
	engine := validator.Direct()
	engine.Handle(&user)
	if err := engine.Unmarshal([]byte(str), &user); err != nil {
		fmt.Println()
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v", user)
}
