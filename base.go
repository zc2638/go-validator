package validator

import (
	"context"
	"net/http"
)

/**
 * Created by zc on 2019-08-12.
 */

type Validation interface {
	Name() string
	SetCondition(...interface{}) error
	Valid(interface{}) error
}

type Hook interface {
	Name() string
	SetCondition(...interface{}) error
	Fire(*VdrEngine) error
}

type Validate interface {
	Register(...Validation)
	SetHook(...Hook)
	SetContext(context.Context) Validate
	SetHttpRequest(*http.Request) Validate
	MakeStruct(interface{}) Validate
	MakeMap(map[string]string) Validate
	MakeSlice(...[]string) Validate
	MakeValue(interface{}, ...string) Validate
	Check() error
}

type VdrEngine struct {
	Name string
	Err error
	Val interface{}
}