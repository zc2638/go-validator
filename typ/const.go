package typ

const (
	NotRequired    = ValidError("string is not required")
	TypeNotSupport = ValidError("type not support")
	TypeNotFound   = ValidError("type not found")
	TypeNotString  = ValidError("type is not string")
	TypeNotInt     = ValidError("type is not int")
	TypeNotInt8    = ValidError("type is not int8")
	TypeNotInt16   = ValidError("type is not int16")
	TypeNotInt32   = ValidError("type is not int32")
	TypeNotInt64   = ValidError("type is not int64")
	TypeNotUint    = ValidError("type is not uint")
	TypeNotUint8   = ValidError("type is not uint8")
	TypeNotUint16  = ValidError("type is not uint16")
	TypeNotUint32  = ValidError("type is not uint32")
	TypeNotUint64  = ValidError("type is not uint64")
	TypeNotFloat32 = ValidError("type is not float32")
	TypeNotFloat64 = ValidError("type is not float64")
	TypeNotBool    = ValidError("type is not bool")
	TypeNotStruct  = ValidError("type is not struct")
	NumberMaxLimit = ValidError("number is over max limit")
	NumberMinLimit = ValidError("number is below min limit")
	RegexpNotMatch = ValidError("string is not match")

	WithoutContext   = ValidError("context is nil")
	SliceFormatError = ValidError("slice format error")
)

const (
	String     = "string"
	Int        = "int"
	Int8       = "int8"
	Int16      = "int16"
	Int32      = "int32"
	Int64      = "int64"
	Uint       = "uint"
	Uint8      = "uint8"
	Uint16     = "uint16"
	Uint32     = "uint32"
	Uint64     = "uint64"
	Float32    = "float32"
	Float64    = "float64"
	Bool       = "bool"
	Complex64  = "complex64"
	Complex128 = "complex128"
)
