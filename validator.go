package validator

import (
	"context"
	"github.com/zc2638/go-validator/typ"
	"net/http"
	"reflect"
	"strings"
)

/**
 * Created by zc on 2019-08-12.
 */

const ValidTagName = "vdr"
const CondMark = ","
const Delimiter = "="

// TODO 是否校验不存在的规则存在时抛错
// TODO 增加一个并发安全的验证器

// 两种验证器类型：1. 父验证器，用于自动创建子验证器 2. 子验证器，用于校验数据
type vdr struct {
	ctx          context.Context
	req          *http.Request
	source       map[string]Validation
	hooks        map[string]Validation
	customSource map[string]Validation
	set          []string
	engines      []*VdrEngine
	err          error
	mark         string
	delimiter    string
}

// 创建一个默认的验证器
func NewVdr() Validate {
	validation := new(vdr)
	validation.register(new(RuleRequired), new(RuleTypes), new(RuleMax), new(RuleMin))
	validation.SetHook(new(HookMsg))
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
	if v.set == nil {
		v.set = make([]string, 0)
	}
	for _, rule := range rules {
		if rule.Name() != "" {
			v.source[rule.Name()] = rule
			v.set = append(v.set, rule.Name())
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
	if v.set == nil {
		v.set = make([]string, 0)
	}
	for _, hook := range hooks {
		if hook.Name() != "" {
			v.hooks[hook.Name()] = hook
			v.set = append(v.set, hook.Name())
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
				// TODO 暂时只支持常用类型，map，slice，struct等后期补充
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
				case reflect.Interface:
					value = fv.Interface()
				case reflect.Map:
				case reflect.Array:
				case reflect.Slice:
				case reflect.Struct:
				case reflect.Ptr:
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
	v.parse(val, exps...)
	return v
}

// TODO 解析条件
func (v *vdr) parse(val interface{}, exps ...string) {

	mark := CondMark
	delimiter := Delimiter
	if v.mark != "" {
		mark = v.mark
	}
	if v.delimiter != "" {
		delimiter = v.delimiter
	}

	expSet := strings.Split(strings.Join(exps, mark), mark)
	var cs = make([]*VdrEngine, 0)
	var comb, expKey string
	for _, e := range expSet {
		if e == "" {
			continue
		}

		// 如果匹配上，则exp赋值
		var exp string
		for _, s := range v.set {
			if strings.HasPrefix(e, s) {
				exp = s
				break
			}
		}

		// 如果exp为空，则表示未匹配上，comb做拼接
		if exp == "" {
			comb += mark + e
			continue
		}

		if expKey == "" {
			expKey = exp
		}

		// 如果exp不为空，则表示匹配上了，将comb内容做消化，comb重新赋值
		if comb != "" {
			expVal := strings.TrimPrefix(strings.TrimPrefix(comb, expKey+delimiter), expKey)
			cs = append(cs, &VdrEngine{
				Name:   expKey,
				Params: []interface{}{expVal},
				Val:    val,
			})
			expKey = exp
		}
		comb = e
	}
	if comb != "" {
		for _, s := range v.set {
			if strings.HasPrefix(comb, s) {
				expVal := strings.TrimPrefix(strings.TrimPrefix(comb, s+delimiter), s)
				cs = append(cs, &VdrEngine{
					Name:   expKey,
					Params: []interface{}{expVal},
					Val:    val,
				})
				break
			}
		}
	}
	if v.engines == nil {
		v.engines = make([]*VdrEngine, 0)
	}
	v.engines = append(v.engines, cs...)
}

func (v *vdr) verify() {
	if v.err != nil {
		return
	}
	if v.engines == nil {
		return
	}

	for _, e := range v.engines {
		if v.err == nil {
			if s, ok := v.source[e.Name]; ok {
				if err := s.SetCondition(e.Params...); err != nil {
					v.err = err
					break
				}
				if err := s.Fire(e); err != nil {
					v.err = err
				}
			}
		}
	}

	for _, e := range v.engines {
		if h, ok := v.source[e.Name]; ok {
			if err := h.SetCondition(e.Params...); err != nil {
				v.err = err
				break
			}
			if err := h.Fire(e); err != nil {
				v.err = err
				break
			}
			if e.Err != nil {
				v.err = e.Err
				break
			}
		}
	}
}

// 清空
func (v *vdr) reset() {
	v.engines = nil
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
	source []Validation
	hooks  []Validation
}

func NewPVdr() *pVdr {
	return new(pVdr)
}

func (v *pVdr) Register(rules ...Validation) {
	if rules == nil || len(rules) == 0 {
		return
	}
	if v.source == nil {
		v.source = make([]Validation, 0)
	}
	v.source = append(v.source, rules...)
}

func (v *pVdr) SetHook(hooks ...Validation) {
	if hooks == nil || len(hooks) == 0 {
		return
	}
	if v.hooks == nil {
		v.hooks = make([]Validation, 0)
	}
	v.hooks = append(v.hooks, hooks...)
}

func (v *pVdr) base() Validate {
	vdr := NewVdr()
	vdr.Register(v.source...)
	vdr.SetHook(v.hooks...)
	return vdr
}

func (v *pVdr) Check() error                            { return v.base().Check() }
func (v *pVdr) SetContext(ctx context.Context) Validate { return v.base().SetContext(ctx) }
func (v *pVdr) SetHttpRequest(r *http.Request) Validate { return v.base().SetHttpRequest(r) }
func (v *pVdr) MakeStruct(s interface{}) Validate       { return v.base().MakeStruct(s) }
func (v *pVdr) MakeMap(ms map[string]string) Validate   { return v.base().MakeMap(ms) }
func (v *pVdr) MakeSlice(set ...[]string) Validate      { return v.MakeSlice(set...) }
func (v *pVdr) MakeValue(val interface{}, exps ...string) Validate {
	return v.base().MakeValue(val, exps...)
}

var vde = NewPVdr()

func Register(rules ...Validation)                       { vde.Register(rules...) }
func SetHook(hooks ...Validation)                        { vde.SetHook(hooks...) }
func SetContext(ctx context.Context) Validate            { return vde.SetContext(ctx) }
func SetHttpRequest(r *http.Request) Validate            { return vde.SetHttpRequest(r) }
func MakeStruct(s interface{}) Validate                  { return vde.MakeStruct(s) }
func MakeMap(ms map[string]string) Validate              { return vde.MakeMap(ms) }
func MakeSlice(set ...[]string) Validate                 { return vde.MakeSlice(set...) }
func MakeValue(val interface{}, exps ...string) Validate { return vde.MakeValue(val, exps...) }
func Check() error                                       { return vde.Check() }
