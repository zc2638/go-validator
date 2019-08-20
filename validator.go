package validator

import (
	"context"
	"errors"
	"github.com/zc2638/go-validator/typ"
	"reflect"
	"strings"
	"unicode"
)

/**
 * Created by zc on 2019-08-12.
 */

const ValidTagName = "vdr"
const CondMark = ","
const Delimiter = "="

// 两种验证器类型：1. 父验证器，用于自动创s建子验证器 2. 子验证器，用于校验数据
type vdr struct {
	ctx          context.Context
	source       map[string]Validation
	hooks        map[string]Validation
	customSource map[string]Validation
	sourceSet    []string
	hookSet      []string
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
	validation.mark = CondMark
	validation.delimiter = Delimiter
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
	if v.sourceSet == nil {
		v.sourceSet = make([]string, 0)
	}
	for _, rule := range rules {
		if rule.Name() != "" {
			v.source[rule.Name()] = rule
			v.sourceSet = append(v.sourceSet, rule.Name())
		}
	}
}

// 注册自定义规则
func (v *vdr) Register(rules ...Validation) {
	v.register(rules...)
}

// 设置hook
func (v *vdr) SetHook(hooks ...Validation) {
	if hooks == nil || len(hooks) == 0 {
		return
	}
	if v.hooks == nil {
		v.hooks = make(map[string]Validation)
	}
	if v.hookSet == nil {
		v.hookSet = make([]string, 0)
	}
	for _, hook := range hooks {
		if hook.Name() != "" {
			v.hooks[hook.Name()] = hook
			v.hookSet = append(v.hookSet, hook.Name())
		}
	}
}

// 添加context，用于键校验
func (v *vdr) SetContext(ctx context.Context) Checker {
	v.ctx = ctx
	return v
}

// struct验证
func (v *vdr) CheckStruct(s interface{}) error {
	if v.ctx == nil {
		v.err = typ.WithoutContext
		return v.err
	}

	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		v.err = typ.StructPtrError
		return v.err
	}
	t = t.Elem()

	ts := typ.NewTypeS(t, reflect.ValueOf(s).Elem())
	fieldSet := ts.GetFiledSet()

	for _, field := range fieldSet {
		fieldName := field.Name
		tag := field.Tag.Get("json")
		if tag != "" {
			if idx := strings.Index(tag, ","); idx != -1 {
				fieldName = tag[:idx]
			} else {
				fieldName = tag
			}
		}

		if err := ts.SetValue(field, v.ctx.Value(fieldName)); err != nil {
			v.err = err
			return v.err
		}
	}
	return v.MakeStruct(s).Check()
}

// map验证
// "id": "required,max=20"
func (v *vdr) CheckMap(ms map[string]string) error {
	if v.ctx == nil {
		v.err = typ.WithoutContext
		return v.err
	}
	for k, m := range ms {
		v.MakeValue(v.ctx.Value(k), m)
	}
	return v.Check()
}

// slice验证
// ["id", "required,max=20", "min=10"], ["age", "required", "max=100"]
func (v *vdr) CheckSlice(set ...[]string) error {
	if v.ctx == nil {
		v.err = typ.WithoutContext
		return v.err
	}
	for _, s := range set {
		if len(s) > 0 {
			v.MakeValue(v.ctx.Value(s[0]), s[1:]...)
		}
	}
	return v.Check()
}

// 创建struct值验证
func (v *vdr) MakeStruct(s interface{}) Validate {
	sv := reflect.ValueOf(s)
	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		sv = sv.Elem()
	}
	if sv.Kind() == reflect.Struct {
		v.handleStruct(sv)
	}
	return v
}

// 值处理
func (v *vdr) handleValue(value reflect.Value, part Part) {
	switch value.Type().Kind() {
	case reflect.String:
		part.Value = value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		part.Value = value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		part.Value = value.Uint()
	case reflect.Float32, reflect.Float64:
		part.Value = value.Float()
	case reflect.Bool:
		part.Value = value.Bool()
	case reflect.Struct:
		v.handleStruct(value)
		part.Value = value.Interface()
	case reflect.Map:
		v.handleMap(value)
		part.Value = value.Interface()
	case reflect.Array, reflect.Slice:
		v.handleSlice(value)
		part.Value = value.Interface()
	case reflect.Ptr:
		if !value.IsNil() {
			v.handleValue(value.Elem(), part)
			return
		}
	case reflect.Interface:
		if !value.IsNil() {
			v.handleValue(reflect.ValueOf(value.Interface()), part)
			return
		}
	}
	v.makePart(part)
}

// struct处理
func (v *vdr) handleStruct(value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		// 根据字段名首字母大小写判断是否私有字段
		if unicode.IsLower(rune(field.Name[0])) {
			v.err = errors.New("field " + field.Name + " is unexported")
			break
		}
		v.handleValue(value.Field(i), Part{
			Key: field.Name,
			Tag: field.Tag.Get(ValidTagName),
		})
	}
}

// map处理
func (v *vdr) handleMap(value reflect.Value) {
	for _, key := range value.MapKeys() {
		mv := value.MapIndex(key)
		if mv.Kind() == reflect.Interface && !mv.IsNil() {
			mv = reflect.ValueOf(mv.Interface())
		}
		if mv.Kind() == reflect.Ptr && !mv.IsNil() {
			mv = mv.Elem()
		}
		if mv.Kind() == reflect.Struct {
			v.handleStruct(mv)
		}
	}
}

// slice处理
func (v *vdr) handleSlice(value reflect.Value) {
	for i := 0; i < value.Len(); i++ {
		sv := value.Index(i)
		if sv.Kind() == reflect.Interface && !sv.IsNil() {
			sv = reflect.ValueOf(sv.Interface())
		}
		if sv.Kind() == reflect.Ptr && !sv.IsNil() {
			sv = sv.Elem()
		}
		if sv.Kind() == reflect.Struct {
			v.handleStruct(sv)
		}
	}
}

// 创建值验证
func (v *vdr) MakeValue(val interface{}, exps ...string) Validate {
	v.makePart(Part{
		Value: val,
		Tag: strings.Join(exps, v.mark),
	})
	return v
}

// part处理
func (v *vdr) makePart(part Part) {
	if part.Tag == "" {
		return
	}
	v.tagParse(part)
}

// 解析条件
func (v *vdr) tagParse(part Part) {

	var set = make([]string, 0)
	set = append(set, v.sourceSet...)
	set = append(set, v.hookSet...)

	expSet := strings.Split(part.Tag, v.mark)
	var sourceEs = make([]*Engine, 0)
	var hookEs = make([]*Engine, 0)
	var comb, expKey string
	for _, e := range expSet {
		if e == "" {
			continue
		}

		// 如果匹配上，则exp赋值
		var exp string
		for _, s := range set {
			if strings.HasPrefix(e, s) {
				exp = s
				break
			}
		}

		// 如果exp为空，则表示未匹配上，comb做拼接
		if exp == "" {
			comb += v.mark + e
			continue
		}

		if expKey == "" {
			expKey = exp
		}

		// 如果exp不为空，则表示匹配上了，将comb内容做消化，comb重新赋值
		if comb != "" {
			expVal := strings.TrimPrefix(strings.TrimPrefix(comb, expKey+v.delimiter), expKey)
			engine := &Engine{Name: expKey, Params: []interface{}{expVal}, Part: part}

			if _, ok := v.source[engine.Name]; ok {
				sourceEs = append(sourceEs, engine)
			}
			if _, ok := v.hooks[engine.Name]; ok {
				hookEs = append(hookEs, engine)
			}
			expKey = exp
		}
		comb = e
	}
	if comb != "" {
		for _, s := range set {
			if strings.HasPrefix(comb, s) {
				expVal := strings.TrimPrefix(strings.TrimPrefix(comb, s+v.delimiter), s)
				engine := &Engine{Name: expKey, Params: []interface{}{expVal}, Part: part}

				if _, ok := v.source[engine.Name]; ok {
					sourceEs = append(sourceEs, engine)
				}
				if _, ok := v.hooks[engine.Name]; ok {
					hookEs = append(hookEs, engine)
				}
				break
			}
		}
	}

	if v.engines == nil {
		v.engines = make([]*VdrEngine, 0)
	}
	v.engines = append(v.engines, &VdrEngine{sourceEs, hookEs})
}

func (v *vdr) verify() {
	if v.err != nil {
		return
	}
	if v.engines == nil {
		return
	}

	for _, e := range v.engines {
		for _, re := range e.Rule {
			if s, ok := v.source[re.Name]; ok {
				if err := s.SetCondition(re.Params...); err != nil {
					v.err = err
					break
				}
				if err := s.Fire(re); err != nil {
					re.Err = err
					for _, he := range e.Hook {
						if h, ok := v.hooks[he.Name]; ok {
							he.Err = re.Err
							if err := h.SetCondition(he.Params...); err != nil {
								he.Err = err
							}
							if err := h.Fire(he); err != nil {
								he.Err = err
							}
							if he.Err != nil {
								re.Err = he.Err
								break
							}
						}
					}
					if re.Err != nil {
						v.err = re.Err
						goto END
					}
				}
			}
		}
	}
END:
}

// 清空
func (v *vdr) reset() {
	v.engines = nil
	v.ctx = nil
	v.err = nil
}

// 验证
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

func (v *pVdr) Check() error                           { return v.base().Check() }
func (v *pVdr) SetContext(ctx context.Context) Checker { return v.base().SetContext(ctx) }
func (v *pVdr) MakeStruct(s interface{}) Validate      { return v.base().MakeStruct(s) }
func (v *pVdr) MakeValue(val interface{}, exps ...string) Validate {
	return v.base().MakeValue(val, exps...)
}

var vde = NewPVdr()

func Register(rules ...Validation)                       { vde.Register(rules...) }
func SetHook(hooks ...Validation)                        { vde.SetHook(hooks...) }
func SetContext(ctx context.Context) Checker             { return vde.SetContext(ctx) }
func MakeStruct(s interface{}) Validate                  { return vde.MakeStruct(s) }
func MakeValue(val interface{}, exps ...string) Validate { return vde.MakeValue(val, exps...) }
func Check() error                                       { return vde.Check() }
