package validator

import (
	"context"
	"reflect"
)

/**
 * Created by zc on 2019-08-12.
 */

// TODO 添加指定校验类型，如果不指定则默认支持的类型全部校验
// TODO 增加结构体条件校验方式
type Validation interface {
	Name() string         // name
	Fire(e *Engine) error // exec
}

type Validate interface {
	Register(...Validation)                    // register rule
	SetHook(...Validation)                     // set hook
	SetContext(context.Context) Checker        // set context
	MakeStruct(interface{}) Validate           // parse struct
	MakeValue(interface{}, ...string) Validate // parse value
	Check() error                              // check
}

type Checker interface {
	CheckStruct(interface{}) error    // parse context with struct
	CheckMap(map[string]string) error // parse context with map
	CheckSlice(...[]string) error     // parse context with slice
}

type VdrEngine struct {
	Rule []*Engine
	Hook []*Engine
}

type Engine struct {
	Name      string // rule name
	Condition string // rule condition
	Err       error  // error
	Part      Part
}

type Part struct {
	Key   string        // key
	Value reflect.Value // value
	Tag   string        // tag
}
