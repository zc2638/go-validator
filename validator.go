/**
 * Created by zc on 2020/6/12.
 */
package validator

import (
	"fmt"
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
		vdr.append(name, &fv)
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

// 将结构体处理的字段中slice类型的赋予验证方法
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

// 将结构体处理的字段中基础类型的赋予验证方法
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
	keys := make(map[string]string)
	for k := range current {
		path := buildSlicePath(k)
		keys[path] = k
	}
	return v.checkLoop(current, keys)
}

func (v *validate) checkLoop(data Current, keys map[string]string) ErrorChains {
	var chains ErrorChains
	// 循环校验结构
	// TODO 暂不考虑指针处理
	for _, vv := range v.vs {
		fmt.Printf("%+v \n", vv)
		if vv.sign == SignSlice {
			chains = append(chains, vv.checkLoop(data, keys)...)
			continue
		}
		// 如果有子集且有校验集，校验通过则继续向下校验，否则不校验子集
		if len(vv.vs) > 0 {
			if len(vv.vfs) > 0 {
				path := keys[vv.path]
				var val interface{}
				if path != "" {
					val = data[path]
				}
				var err error
				for _, vf := range vv.vfs {
					if err = vf(val); err != nil {
						break
					}
				}
				if err != nil {
					chains = append(chains, Error{
						path: vv.path,
						e:    err,
					})
					continue
				}
			}
			// 如果类型为切片，入参无值跳过，如 arr: []
			if vv.name == SignSlice {
				path := keys[vv.path]
				val, ok := data[path]
				if !ok {
					continue
				}
				// slice长度为空，跳过
				if reflect.ValueOf(val).Len() == 0 {
					continue
				}
			}
			chains = append(chains, vv.checkLoop(data, keys)...)
		} else {
			if len(vv.vfs) == 0 {
				continue
			}
			path := keys[vv.path]
			val := data[path]
			for _, vf := range vv.vfs {
				if err := vf(val); err != nil {
					chains = append(chains, Error{
						path: vv.path,
						e:    err,
					})
					break
				}
			}
		}
	}
	return chains
}
