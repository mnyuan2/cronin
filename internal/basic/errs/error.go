package errs

import (
	"fmt"
	"runtime"
)

type ErrBase interface {
	error
	Code() string
	Path() string // 调用错误路径
	ReloadPath()
}

type ErrorOption interface {
	parse(e *Error)
}

// 错误码
type Code string

func (c Code) parse(e *Error) {
	e.code = c
}

const (
	SysError = Code("999999") // 系统错误
)

// 错误描述
type Desc string

func (c Desc) parse(e *Error) {
	e.desc = string(c)
}

type Error struct {
	code Code   // 错误码
	err  error  // 原始错误
	desc string // 自定义描述
	path string // 调用者路径
}

// New 返回错误
func New(err error, ops ...any) ErrBase {
	_, file, line, _ := runtime.Caller(1)

	e := &Error{
		err:  err,
		path: fmt.Sprintf("%s %v", file, line),
	}

	for _, op := range ops {
		switch op.(type) {
		case Code:
			op.(Code).parse(e)
		case string, Desc:
			e.desc = op.(string)
		}
	}
	return e
}

// ReloadPath 重新加载错误的触发路径
func (e *Error) ReloadPath() {
	_, file, line, _ := runtime.Caller(1)
	e.path = fmt.Sprintf("%s %v", file, line)
}

// Code 错误码
func (e *Error) Code() string {
	return string(e.code)
}

// Desc 自定义错误描述
func (e *Error) Desc() string {
	return e.desc
}

// Path 错误触发点
func (e *Error) Path() string {
	return e.path
}

// Error 原始错误描述
func (e *Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}
