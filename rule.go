package validator

import (
	"errors"
)

/**
 * Created by zc on 2019-08-13.
 */

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
