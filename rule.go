package validator

import (
	"errors"
	"reflect"
	"strconv"
)

/**
 * Created by zc on 2019-08-13.
 */

func RuleType(kind reflect.Kind) ValidateFunc {
	return func(val interface{}) error {
		if val == nil {
			return nil
		}
		var message string
		switch val.(type) {
		case bool:
			if kind != reflect.Bool {
				message = "type exception: bool"
			}
		case string:
			if kind != reflect.String {
				message = "type exception: string"
			}
		case float64:
			var err error
			str := strconv.FormatFloat(val.(float64), 'f', -1, 64)
			switch kind {
			case reflect.Uint:
				_, err = strconv.ParseUint(str, 10, strconv.IntSize)
			case reflect.Uint8:
				_, err = strconv.ParseUint(str, 10, 8)
			case reflect.Uint16:
				_, err = strconv.ParseUint(str, 10, 16)
			case reflect.Uint32:
				_, err = strconv.ParseUint(str, 10, 32)
			case reflect.Uint64:
				_, err = strconv.ParseUint(str, 10, 64)
			case reflect.Int:
				_, err = strconv.Atoi(str)
			case reflect.Int8:
				_, err = strconv.ParseInt(str, 10, 8)
			case reflect.Int16:
				_, err = strconv.ParseInt(str, 10, 16)
			case reflect.Int32:
				_, err = strconv.ParseInt(str, 10, 32)
			case reflect.Int64:
				_, err = strconv.ParseInt(str, 10, 64)
			case reflect.Float32:
				_, err = strconv.ParseFloat(str, 32)
			case reflect.Float64:
				return nil
			default:
				message = "type exception: " + kind.String() + " unsupport"
			}
			if err != nil {
				message = "type exception: " + kind.String()
			}
		default:
			message = "type exception"
		}
		if message != "" {
			return errors.New(message)
		}
		return nil
	}
}

func autoMountRuleType(kind reflect.Kind) ValidateFunc {
	if kind > 0 && kind < 12 || kind == reflect.Float32 || kind == reflect.Float64 || kind == reflect.String {
		return RuleType(kind)
	}
	return nil
}

func RuleRequired() ValidateFunc {
	return RuleRequiredWithMessage("not required")
}

func RuleRequiredWithMessage(message string) ValidateFunc {
	return func(val interface{}) error {
		if val == nil {
			return errors.New(message)
		}
		return nil
	}
}
