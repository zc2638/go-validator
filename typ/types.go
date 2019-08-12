package typ

import (
	"reflect"
	"strconv"
)

type ValidError string
func (e ValidError) Error() string { return string(e) }

type typeC struct {
	val  interface{}
	kind reflect.Kind
	typ  string
}

func NewTypeC(val interface{}, typ string) *typeC {
	rt := reflect.TypeOf(val)
	return &typeC{val: val, typ: typ, kind: rt.Kind()}
}

// 检查类型匹配
func (t *typeC) CheckKind() bool {
	switch t.kind {
	case reflect.String:
		return t.typ == String
	case reflect.Int:
		return t.typ == Int
	case reflect.Int8:
		return t.typ == Int8
	case reflect.Int16:
		return t.typ == Int16
	case reflect.Int32:
		return t.typ == Int32
	case reflect.Int64:
		return t.typ == Int64
	case reflect.Uint:
		return t.typ == Uint
	case reflect.Uint8:
		return t.typ == Uint8
	case reflect.Uint16:
		return t.typ == Uint16
	case reflect.Uint32:
		return t.typ == Uint32
	case reflect.Uint64:
		return t.typ == Uint64
	case reflect.Float32:
		return t.typ == Float32
	case reflect.Float64:
		return t.typ == Float64
	default:
		return false
	}
}

// 检查基础类型
func (t *typeC) Convert() (interface{}, error) {

	var str string
	switch t.kind {
	case reflect.String:
		str = t.val.(string)
	case reflect.Int:
		str = strconv.Itoa(t.val.(int))
		return t.stringToType(str)
	case reflect.Int8:
		str = strconv.Itoa(int(t.val.(int8)))
	case reflect.Int16:
		str = strconv.Itoa(int(t.val.(int16)))
	case reflect.Int32:
		str = strconv.Itoa(int(t.val.(int32)))
	case reflect.Int64:
		str = strconv.Itoa(int(t.val.(int64)))
	case reflect.Uint:
		str = strconv.Itoa(int(t.val.(uint)))
	case reflect.Uint8:
		str = strconv.Itoa(int(t.val.(uint8)))
	case reflect.Uint16:
		str = strconv.Itoa(int(t.val.(uint16)))
	case reflect.Uint32:
		str = strconv.Itoa(int(t.val.(uint32)))
	case reflect.Uint64:
		str = strconv.Itoa(int(t.val.(uint64)))
	case reflect.Float32:
		str = strconv.FormatFloat(float64(t.val.(float32)), 'f', -1, 32)
	case reflect.Float64:
		str = strconv.FormatFloat(t.val.(float64), 'f', -1, 64)
	default:
		return nil, TypeNotSupport
	}
	return t.stringToType(str)
}

func (t *typeC) stringToType(str string) (interface{}, error) {

	var res interface{}
	var err error = TypeNotFound

	switch t.typ {
	case String:
		return str, nil
	case Int:
		res, err = strconv.Atoi(str)
		if err != nil {
			err = TypeNotInt
		}
	case Int8:
		r, e := strconv.ParseInt(str, 10, 8)
		if e != nil {
			err = TypeNotInt8
		} else {
			res = int8(r)
		}
	case Int16:
		r, e := strconv.ParseInt(str, 10, 16)
		if e != nil {
			err = TypeNotInt16
		} else {
			res = int16(r)
		}
	case Int32:
		r, e := strconv.ParseInt(str, 10, 32)
		if e != nil {
			err = TypeNotInt32
		} else {
			res = int32(r)
		}
	case Int64:
		r, e := strconv.ParseInt(str, 10, 64)
		if e != nil {
			err = TypeNotInt64
		} else {
			res = int64(r)
		}
	case Uint:
		r, e := strconv.ParseUint(str, 10, 0)
		if e != nil {
			err = TypeNotUint
		} else {
			res = uint(r)
		}
	case Uint8:
		r, e := strconv.ParseUint(str, 10, 8)
		if e != nil {
			err = TypeNotUint8
		} else {
			res = uint8(r)
		}
	case Uint16:
		r, e := strconv.ParseUint(str, 10, 16)
		if e != nil {
			err = TypeNotUint16
		} else {
			res = uint16(r)
		}
	case Uint32:
		r, e := strconv.ParseUint(str, 10, 32)
		if e != nil {
			err = TypeNotUint32
		} else {
			res = uint32(r)
		}
	case Uint64:
		res, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			err = TypeNotUint64
		}
	case Float32:
		r, e := strconv.ParseFloat(str, 32)
		if e != nil {
			err = TypeNotFloat32
		} else {
			res = float32(r)
		}
	case Float64:
		res, err = strconv.ParseFloat(str, 64)
		if err != nil {
			err = TypeNotFloat64
		}
	}
	return res, err
}
