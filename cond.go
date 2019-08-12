package validator

import (
	"strings"
)

/**
 * Created by zc on 2019-08-12.
 */

type Cond struct {
	key   string      // 键
	value interface{} // 值
	cs    []Condition // 验证条件
	err   error       // 错误
}

// 条件解析结构
// TODO 正则分割的时候会出问题，后期支持正则
const CondMark = ","

type Condition struct {
	Name  string
	value []interface{}
}

func newCond(value interface{}, err error, exps ...string) *Cond {
	if len(exps) == 0 {
		return &Cond{value: value}
	}

	c := make([]Condition, 0)
	expSet := strings.Split(strings.Join(exps, CondMark), CondMark)
	for _, e := range expSet {
		if e == "" {
			continue
		}
		es := strings.Split(e, "=")
		value := make([]interface{}, 0)
		if len(es) > 1 {
			for _, v := range es[1:] {
				value = append(value, v)
			}
		}
		c = append(c, Condition{
			Name:  es[0],
			value: value,
		})
	}

	return &Cond{value: value, cs: c, err: err}
}
