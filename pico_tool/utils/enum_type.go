package utils

type EnumType struct {
	code int
	desc string
}

func NewEnumType(code int, desc string) EnumType {
	return EnumType{
		code: code,
		desc: desc,
	}
}

func (e EnumType) Code() int {
	return e.code
}

func (e EnumType) Desc() string {
	return e.desc
}

func (e EnumType) Value() int {
	return e.code
}

type StatusEnumType struct {
	EnumType
	err error
}

func NewStatusEnumType(code int, desc string, err error) StatusEnumType {
	return StatusEnumType{
		EnumType: NewEnumType(code, desc),
		err:      err,
	}
}

func (e StatusEnumType) Error() error {
	return e.err
}

func (e StatusEnumType) WithError(err error) StatusEnumType {
	e.err = err
	return e
}
