/**
 * Created by zc on 2020/6/12.
 */
package validator

import (
	"reflect"
)

type validate struct {
	tag   string         // 解析使用的tag名称
	sign  string         // 类型标识
	name  string         // 字段名
	path  string         // 校验路径
	value *reflect.Value // 反射值，用于匹配结构体字段
	vfs   []ValidateFunc // 校验方法集
	vs    []*validate    // 子集
}

func newValidate() *validate {
	return &validate{
		tag: TagJSON,
	}
}

func (v *validate) append(name string, rv *reflect.Value, vfs ...ValidateFunc) {
	vdr := newValidate()
	vdr.name = name
	if v.sign != "" {
		vdr.path = buildPath(v.path, v.sign, name)
	} else {
		vdr.path = buildPath(v.path, name)
	}
	vdr.value = rv
	vdr.vfs = vfs
	v.vs = append(v.vs, vdr)
}

func (v *validate) update(rv reflect.Value, vfs ...ValidateFunc) {
	if len(vfs) == 0 {
		return
	}
	for _, vp := range v.vs {
		if *vp.value == rv {
			vp.vfs = vfs
			return
		}
	}
}

func (v *validate) find(rv reflect.Value) *validate {
	for _, vv := range v.vs {
		if *vv.value == rv {
			return vv
		}
	}
	vdr := newValidate()
	if v.sign == SignSlice {
		vdr.name = SignSlice
		vdr.path = buildPath(v.path, vdr.name)
	} else {
		vdr.name = rv.Type().Name()
	}
	v.vs = append(v.vs, vdr)
	return vdr
}

// 解构struct
func (v *validate) decompose(rv reflect.Value, vdr *validate) {
	rt := rv.Type()
	rtFieldNum := rt.NumField()
	for i := 0; i < rtFieldNum; i++ {
		ft := rt.Field(i)
		var name string
		if tagValue, ok := ft.Tag.Lookup(vdr.tag); ok {
			name = tagValue
		} else {
			name = camelToUnderline(ft.Name)
		}
		fv := rv.Field(i)
		vfType := autoMountRuleType(ft.Type.Kind())
		if vfType == nil {
			vdr.append(name, &fv)
		} else {
			vdr.append(name, &fv, vfType)
		}
	}
}

// 处理结构体所有字段
func (v *validate) Make(s Handler, vfs ...ValidateFunc) {
	rv := Value(s)
	if !rv.CanAddr() {
		return
	}
	if rv.Kind() != reflect.Struct {
		return
	}
	v.update(rv, vfs...)
	vdr := v.find(rv)
	vdr.decompose(rv, vdr)
	s.Validate(vdr)
}

// 处理slice
func (v *validate) MakeSlice(val interface{}, f HandlerFunc, vfs ...ValidateFunc) {
	rv := Value(val)
	if rv.Type().Kind() != reflect.Slice {
		return
	}
	v.update(rv, vfs...)
	vdr := v.find(rv)
	vdr.sign = SignSlice
	if f != nil {
		f(vdr)
	}
}

// 处理基础值
func (v *validate) MakeValue(val interface{}, vfs ...ValidateFunc) {
	rv := Value(val)
	if !rv.CanAddr() {
		return
	}
	v.update(rv, vfs...)
}

func (v *validate) MakeField(name string, vfs ...ValidateFunc) {
	v.append(name, nil, vfs...)
}

func (v *validate) Check(current Current) error {
	// 由于切片入参是多组，所以key对应源path组
	keys := make(map[string][]string)
	for k := range current {
		path := buildSlicePath(k)
		keys[path] = append(keys[path], k)
	}
	if err := v.checkLoop(current, keys, ""); err != nil {
		return err
	}
	return nil
}

func (v *validate) verify(val interface{}, vfs []ValidateFunc) error {
	for _, vf := range vfs {
		if err := vf(val); err != nil {
			return err
		}
	}
	return nil
}

func (v *validate) checkLoop(data Current, keys map[string][]string, pre string) ErrorChains {
	var chains ErrorChains
	// 循环校验结构
	// TODO 暂不考虑指针处理
	for _, vv := range v.vs {
		if vv.path == "" {
			chains = append(chains, vv.checkLoop(data, keys, "")...)
			continue
		}
		// 拼接path
		path := buildPath(pre, vv.name)
		// 如果类型为切片，入参无值跳过，如 arr: []
		if vv.name == SignSlice {
			transPaths := keys[vv.path]
			if len(transPaths) > 0 {
				for _, tp := range transPaths {
					chains = append(chains, vv.checkLoop(data, keys, tp)...)
				}
			}
			continue
		}

		// 如果有校验集，校验通过则继续向下校验，否则不校验子集
		if len(vv.vfs) > 0 {
			transPaths := keys[vv.path]
			var currentPath string
			for _, tp := range transPaths {
				if tp == path {
					currentPath = tp
					break
				}
			}
			// 匹配到正常校验，未匹配到按nil校验
			var val interface{}
			if currentPath != "" {
				val = data[currentPath]
			}
			var err error
			for _, vf := range vv.vfs {
				if err = vf(val); err != nil {
					break
				}
			}
			if err != nil {
				chains = append(chains, Error{
					path: path,
					err:  err,
				})
				continue
			}
		}
		if len(vv.vs) > 0 {
			chains = append(chains, vv.checkLoop(data, keys, path)...)
		}
	}
	return chains
}

func (v *validate) checkLoopDirect(data Current, keys map[string][]string, pre string) *Error {
	for _, vv := range v.vs {
		if vv.path == "" {
			if err := vv.checkLoopDirect(data, keys, ""); err != nil {
				return err
			}
			continue
		}
		path := buildPath(pre, vv.name)
		if vv.name == SignSlice {
			transPaths := keys[vv.path]
			if len(transPaths) > 0 {
				for _, tp := range transPaths {
					if err := vv.checkLoopDirect(data, keys, tp); err != nil {
						return err
					}
				}
			}
			continue
		}
		if len(vv.vfs) > 0 {
			transPaths := keys[vv.path]
			var currentPath string
			for _, tp := range transPaths {
				if tp == path {
					currentPath = tp
					break
				}
			}
			var val interface{}
			if currentPath != "" {
				val = data[currentPath]
			}
			for _, vf := range vv.vfs {
				if err := vf(val); err != nil {
					return &Error{path: path, err: err}
				}
			}
		}
		if len(vv.vs) > 0 {
			if err := vv.checkLoopDirect(data, keys, path); err != nil {
				return err
			}
		}
	}
	return nil
}

type validateDirect struct {
	validate
}

func newValidateDirect() *validateDirect {
	return &validateDirect{
		validate: validate{
			tag: TagJSON,
		},
	}
}

func (v *validateDirect) Check(current Current) error {
	keys := make(map[string][]string)
	for k := range current {
		path := buildSlicePath(k)
		keys[path] = append(keys[path], k)
	}
	if err := v.checkLoopDirect(current, keys, ""); err != nil {
		return err
	}
	return nil
}
