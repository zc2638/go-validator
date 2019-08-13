package validator

import (
	"context"
	"github.com/zc2638/go-validator/typ"
	"net/http"
	"reflect"
)

/**
 * Created by zc on 2019-08-12.
 */

const ValidTagName = "vdr"

// TODO 是否校验不存在的规则存在时抛错
// TODO 增加一个并发安全的验证器

// 两种验证器类型：1. 父验证器，用于自动创建子验证器 2. 子验证器，用于校验数据
type vdr struct {
	ctx          context.Context
	req          *http.Request
	source       map[string]Validation
	hooks        map[string]Validation
	customSource map[string]Validation
	cds          []*Cond
	err          error
}

// 创建一个默认的验证器
func NewVdr() Validate {

	required := new(RuleRequired)
	types := new(RuleTypes)
	max := new(RuleMax)
	min := new(RuleMin)

	validation := new(vdr)
	validation.source = map[string]Validation{
		required.Name(): required,
		types.Name():    types,
		max.Name():      max,
		min.Name():      min,
	}

	msg := new(HookMsg)
	validation.hooks = map[string]Validation{
		msg.Name(): msg,
	}
	return validation
}

// 注册初始规则
func (v *vdr) register(rules ...Validation) {
	if rules == nil || len(rules) == 0 {
		return
	}
	if v.source == nil {
		v.source = make(map[string]Validation)
	}
	for _, rule := range rules {
		if rule.Name() != "" {
			v.source[rule.Name()] = rule
		}
	}
}

// 注册自定义规则
func (v *vdr) Register(rules ...Validation) {
	if rules == nil || len(rules) == 0 {
		return
	}
	if v.customSource == nil {
		v.customSource = make(map[string]Validation)
	}
	for _, rule := range rules {
		if rule.Name() != "" {
			v.customSource[rule.Name()] = rule
		}
	}
}

func (v *vdr) SetHook(hooks ...Validation) {
	if hooks == nil || len(hooks) == 0 {
		return
	}
	if v.hooks == nil {
		v.hooks = make(map[string]Validation)
	}
	for _, hook := range hooks {
		if hook.Name() != "" {
			v.hooks[hook.Name()] = hook
		}
	}
}

// 添加context，用于键校验
func (v *vdr) SetContext(context context.Context) Validate {
	v.ctx = context
	return v
}

// 添加*http.Request,用于校验
func (v *vdr) SetHttpRequest(r *http.Request) Validate {
	v.req = r
	return v
}

// 创建struct验证
func (v *vdr) MakeStruct(s interface{}) Validate {

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		t := val.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fv := val.Field(i)

			if tag, ok := field.Tag.Lookup(ValidTagName); ok {
				// TODO 暂时只支持常用类型，byte，map，slice，struct等后期补充
				var value interface{}
				switch field.Type.Kind() {
				case reflect.String:
					value = fv.String()
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value = fv.Int()
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					value = fv.Uint()
				case reflect.Float32, reflect.Float64:
					value = fv.Float()
				case reflect.Bool:
					value = fv.Bool()
				}
				if value == nil {
					continue
				}
				v.MakeValue(value, tag)
			}
		}
	}
	return v
}

// 创建map验证
// "id": "required,max=20"
func (v *vdr) MakeMap(ms map[string]string) Validate {
	if v.ctx == nil {
		v.err = typ.WithoutContext
		return v
	}
	for k, m := range ms {
		v.MakeValue(v.ctx.Value(k), m)
	}
	return v
}

// 创建slice验证
// ["id", "required,max=20", "min=10"], ["age", "required", "max=100"]
func (v *vdr) MakeSlice(set ...[]string) Validate {
	if v.ctx == nil {
		v.err = typ.WithoutContext
		return v
	}
	for _, s := range set {
		if len(s) > 0 {
			v.MakeValue(v.ctx.Value(s[0]), s[1:]...)
		}
	}
	return v
}

// 创建值验证
func (v *vdr) MakeValue(val interface{}, exps ...string) Validate {
	v.cds = append(v.cds, newCond(val, nil, exps...))
	return v
}

// 校验
func (v *vdr) verify() {
	if v.err != nil {
		return
	}
	if v.cds == nil {
		return
	}

	var hookCds = make([]*Cond, 0)
	for _, cond := range v.cds {

		if v.err != nil {
			break
		}
		var hookCond = &Cond{
			key:   cond.key,
			value: cond.value,
		}
		for _, condition := range cond.cs {

			engine := &VdrEngine{
				Name: condition.Name,
				Err:  v.err,
				Val:  cond.value,
			}

			if v.err == nil {
				if s, ok := v.source[condition.Name]; ok {
					if err := s.SetCondition(condition.value...); err != nil {
						v.err = err
						break
					}
					if err := s.Fire(engine); err != nil {
						if cond.err != nil {
							v.err = cond.err
						} else {
							v.err = err
						}
					}
				}
			}
			if v.err == nil {
				if sc, ok := v.customSource[condition.Name]; ok {
					if err := sc.SetCondition(condition.value...); err != nil {
						v.err = err
						break
					}
					if err := sc.Fire(engine); err != nil {
						if cond.err != nil {
							v.err = cond.err
						} else {
							v.err = err
						}
					}
				}
			}

			_, ok := v.hooks[condition.Name]
			if ok {
				if hookCond.cs == nil {
					hookCond.cs = make([]Condition, 0)
				}
				hookCond.cs = append(hookCond.cs, Condition{
					Name:  condition.Name,
					value: condition.value,
				})
			}
		}

		if hookCond.cs != nil {
			hookCds = append(hookCds, hookCond)
		}
	}

	if len(hookCds) == 0 {
		return
	}
	for _, cond := range hookCds {
		for _, condition := range cond.cs {
			h := v.hooks[condition.Name]
			if err := h.SetCondition(condition.value...); err != nil {
				v.err = err
				break
			}
			engine := &VdrEngine{
				Name: condition.Name,
				Err:  v.err,
				Val:  cond.value,
			}

			if err := h.Fire(engine); err != nil {
				v.err = err
				break
			}
			v.setEngine(engine)
		}
	}
}

func (v *vdr) setEngine(e *VdrEngine) {
	v.err = e.Err
}

// 清空
func (v *vdr) reset() {
	v.cds = nil
	v.ctx = nil
	v.req = nil
	v.err = nil
}

// 验证结果
func (v *vdr) Check() error {
	v.verify()
	err := v.err
	v.reset()
	return err
}

type pVdr struct {
	source       map[string]Validation
	customSource map[string]Validation
}

func NewPVdr() *pVdr {
	return new(pVdr)
}

func (v *pVdr) register(rules ...Validation) {
	if rules == nil || len(rules) == 0 {
		return
	}
	if v.source == nil {
		v.source = make(map[string]Validation)
	}
	for _, rule := range rules {
		if rule.Name() != "" {
			v.source[rule.Name()] = rule
		}
	}
}

func (v *pVdr) Register(rules ...Validation) {
	if rules == nil || len(rules) == 0 {
		return
	}
	if v.source == nil {
		v.source = make(map[string]Validation)
	}
	for _, rule := range rules {
		if rule.Name() != "" {
			v.customSource[rule.Name()] = rule
		}
	}
}

func (v *pVdr) Check() error {
	return v.base().Check()
}

func (v *pVdr) SetContext(ctx context.Context) Validate {
	return v.base().SetContext(ctx)
}

func (v *pVdr) SetHttpRequest(r *http.Request) Validate {
	return v.base().SetHttpRequest(r)
}

func (v *pVdr) MakeStruct(s interface{}) Validate {
	return v.base().MakeStruct(s)
}

func (v *pVdr) MakeMap(ms map[string]string) Validate {
	return v.base().MakeMap(ms)
}

func (v *pVdr) MakeSlice(set ...[]string) Validate {
	return v.MakeSlice(set...)
}

func (v *pVdr) MakeValue(val interface{}, exps ...string) Validate {
	return v.base().MakeValue(val, exps...)
}

func (v *pVdr) base() Validate {
	return &vdr{
		source:       v.source,
		customSource: v.customSource,
	}
}

var vde = NewPVdr()

func Register(rules ...Validation) { vde.Register(rules...) }

func SetContext(ctx context.Context) Validate { return vde.SetContext(ctx) }

func SetHttpRequest(r *http.Request) Validate { return vde.SetHttpRequest(r) }

func MakeStruct(s interface{}) Validate { return vde.MakeStruct(s) }

func MakeMap(ms map[string]string) Validate { return vde.MakeMap(ms) }

func MakeSlice(set ...[]string) Validate { return vde.MakeSlice(set...) }

func MakeValue(val interface{}, exps ...string) Validate { return vde.MakeValue(val, exps...) }

func Check() error { return vde.Check() }
