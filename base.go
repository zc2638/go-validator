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
	Fire(*VdrEngine) error             // exec
}

type Validate interface {
	Register(...Validation)                    // register rule
	SetHook(...Validation)                     // set hook
	SetContext(context.Context) Validate       // set context
	MakeStruct(interface{}) Validate           // parse context with struct
	MakeMap(map[string]string) Validate        // parse context with map
	MakeSlice(...[]string) Validate            // parse context with slice
	MakeStructValue(interface{}) Validate      // parse struct
	MakeValue(interface{}, ...string) Validate // parse value
	Check() error                              // check
}

type VdrEngine struct {
	Name   string        // rule name
	Params []interface{} // rule params
	Err    error         // error
	Key    string        // key
	Val    interface{}   // value
}
