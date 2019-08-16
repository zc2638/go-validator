package typ

import (
	"reflect"
	"strconv"
)

type ValidError string

func (e ValidError) Error() string { return string(e) }

func ChangeTypeToKind(typ string) (kind reflect.Kind) {
	switch typ {
	case String:
		kind = reflect.String
	case Int:
		kind = reflect.Int
	case Int8:
		kind = reflect.Int8
	case Int16:
		kind = reflect.Int16
	case Int32:
		kind = reflect.Int32
	case Int64:
		kind = reflect.Int64
	case Uint:
		kind = reflect.Uint
	case Uint8:
		kind = reflect.Uint8
	case Uint16:
		kind = reflect.Uint16
	case Uint32:
		kind = reflect.Uint32
	case Uint64:
		kind = reflect.Uint64
	case Bool:
		kind = reflect.Bool
	case Complex64:
		kind = reflect.Complex64
	case Complex128:
		kind = reflect.Complex128
	}
	return
}

type typeC struct {
	val interface{}
	k   reflect.Kind
	ck  reflect.Kind
}

func NewTypeC(val interface{}, kind reflect.Kind) (tc *typeC) {
	tc = &typeC{ck: kind}
	if val != nil {
		tc.val = val
		tc.k = reflect.TypeOf(val).Kind()
	}
	return
}

// 检查类型匹配
func (t *typeC) CheckKind() bool {
	return t.k == t.ck
}

// 检查基础类型
func (t *typeC) Convert() (interface{}, error) {

	var str string
	switch t.k {
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
	case reflect.Bool:
		str = strconv.FormatBool(t.val.(bool))
	default:
		return nil, TypeNotSupport
	}
	return t.stringToType(str)
}

func (t *typeC) stringToType(str string) (interface{}, error) {

	var res interface{}
	var err error = TypeNotFound

	switch t.ck {
	case reflect.Interface:
		return str, nil
	case reflect.String:
		return str, nil
	case reflect.Int:
		res, err = strconv.Atoi(str)
		if err != nil {
			err = TypeNotInt
		}
	case reflect.Int8:
		r, e := strconv.ParseInt(str, 10, 8)
		if e != nil {
			err = TypeNotInt8
		} else {
			res = int8(r)
		}
	case reflect.Int16:
		r, e := strconv.ParseInt(str, 10, 16)
		if e != nil {
			err = TypeNotInt16
		} else {
			res = int16(r)
		}
	case reflect.Int32:
		r, e := strconv.ParseInt(str, 10, 32)
		if e != nil {
			err = TypeNotInt32
		} else {
			res = int32(r)
		}
	case reflect.Int64:
		r, e := strconv.ParseInt(str, 10, 64)
		if e != nil {
			err = TypeNotInt64
		} else {
			res = int64(r)
		}
	case reflect.Uint:
		r, e := strconv.ParseUint(str, 10, 0)
		if e != nil {
			err = TypeNotUint
		} else {
			res = uint(r)
		}
	case reflect.Uint8:
		r, e := strconv.ParseUint(str, 10, 8)
		if e != nil {
			err = TypeNotUint8
		} else {
			res = uint8(r)
		}
	case reflect.Uint16:
		r, e := strconv.ParseUint(str, 10, 16)
		if e != nil {
			err = TypeNotUint16
		} else {
			res = uint16(r)
		}
	case reflect.Uint32:
		r, e := strconv.ParseUint(str, 10, 32)
		if e != nil {
			err = TypeNotUint32
		} else {
			res = uint32(r)
		}
	case reflect.Uint64:
		res, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			err = TypeNotUint64
		}
	case reflect.Float32:
		r, e := strconv.ParseFloat(str, 32)
		if e != nil {
			err = TypeNotFloat32
		} else {
			res = float32(r)
		}
	case reflect.Float64:
		res, err = strconv.ParseFloat(str, 64)
		if err != nil {
			err = TypeNotFloat64
		}
	case reflect.Bool:
		res, err = strconv.ParseBool(str)
		if err != nil {
			err = TypeNotBool
		}
	}
	return res, err
}
