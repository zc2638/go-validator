package validator

import (
	"context"
)

/**
 * Created by zc on 2019-08-12.
 */

type Validation interface {
	Name() string                      // name
	SetCondition(...interface{}) error // set condition
	Fire(*Engine) error                // exec
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
	Name   string
	Params []interface{}
	Err    error
	Key    string
	Val    interface{}
}
