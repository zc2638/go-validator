package validator

import (
	"errors"
)

/**
 * Created by zc on 2019-08-13.
 */

type HookMsg struct{}

func (h *HookMsg) Name() string { return "msg" }
func (h *HookMsg) Fire(e *Engine) error {
	if e.Err != nil {
		e.Err = errors.New(e.Condition)
	}
	return nil
}
