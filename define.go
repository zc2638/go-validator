package validator

/**
 * Created by zc on 2019-08-12.
 */

const TagJSON = "json"

const SignSlice = "0slice"

// 验证定义 - 用户结构实现
type Handler interface {
	Validate(Validation)
}

type HandlerFunc func(v Validation)

func (f HandlerFunc) Validate(v Validation) { f(v) }

type Validation interface {
	Validator
	Checker
}

type Checker interface {
	Check(current Current) error
}

// 结构解析定义
type Validator interface {
	Make(handler Handler, vfs ...ValidateFunc)                   // parse struct
	MakeSlice(s interface{}, f HandlerFunc, vfs ...ValidateFunc) // parse slice
	MakeValue(val interface{}, vfs ...ValidateFunc)              // parse value
	MakeField(name string, vfs ...ValidateFunc)                  // parse field
}

// 校验方法
type ValidateFunc func(val interface{}) error

// 接收值定义
type Current map[string]interface{}

// 解析器
type Formatter func(data []byte) (Current, error)

// 构造器
type Cover func(data []byte, s interface{}) error

// 校验合并
func Register(vfs ...ValidateFunc) ValidateFunc {
	return func(val interface{}) error {
		for _, vf := range vfs {
			if err := vf(val); err != nil {
				return err
			}
		}
		return nil
	}
}

type Engine struct {
	formatter  Formatter  // 入参解析器
	cover      Cover      // 出参解析器
	validation Validation // 处理器
}

func Default() *Engine {
	e := &Engine{}
	e.formatter = JSONFormatter()
	e.cover = JSONCover()
	e.validation = newValidate()
	return e
}

func Direct() *Engine {
	e := &Engine{}
	e.formatter = JSONFormatter()
	e.cover = JSONCover()
	e.validation = newValidateDirect()
	return e
}

func (e *Engine) SetFormatter(formatter Formatter) {
	if formatter != nil {
		e.formatter = formatter
	}
}

func (e *Engine) SetCover(cover Cover) {
	if cover != nil {
		e.cover = cover
	}
}

func (e *Engine) Handle(s Handler, vfs ...ValidateFunc) Validator {
	e.validation.Make(s, vfs...)
	return e.validation
}

func (e *Engine) HandleSlice(s interface{}, f HandlerFunc, vfs ...ValidateFunc) Validator {
	e.validation.MakeSlice(s, f, vfs...)
	return e.validation
}

func (e *Engine) Check(data []byte) error {
	current, err := e.formatter(data)
	if err != nil {
		return err
	}
	return e.validation.Check(current)
}

func (e *Engine) Unmarshal(data []byte, v interface{}) error {
	if err := e.Check(data); err != nil {
		return err
	}
	return e.cover(data, v)
}
