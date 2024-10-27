package db

import (
	"reflect"
)

const (
	DriverMysql  = "mysql"
	DriverSqlite = "sqlite"
)

type whereOptions struct {
	// 是否必须, 如果为true ,忽略空值
	required bool
	// 默认的空值
	emptyVal interface{}
	// 链接符
	withSpace string
}

func (w whereOptions) isZero(value interface{}) bool {
	valueOf := reflect.ValueOf(value)
	if w.emptyVal == nil {
		return valueOf.IsZero()
	} else {
		kind := valueOf.Kind()

		if kind == reflect.String {
			return value == w.emptyVal
		}
		if kind >= reflect.Int && kind <= reflect.Uint64 {
			return valueOf.Int() == reflect.ValueOf(w.emptyVal).Int()
		}
		if kind >= reflect.Float32 && kind <= reflect.Float64 {
			return valueOf.Float() == reflect.ValueOf(w.emptyVal).Float()
		}

		panic("不支持的数据类型")
	}
}

func ApplyOptions(opts ...Option) *whereOptions {
	whereOpts := &whereOptions{
		withSpace: AndWithSpace,
	}

	for _, item := range opts {
		item.apply(whereOpts)
	}

	return whereOpts
}

type Option interface {
	apply(*whereOptions)
}

// 标记当前查询字段是必须的查询条件, 即使传入的值是空值, 也会查询该条件
type requiredOption int

func (requiredOption) apply(opt *whereOptions) {
	opt.required = true
}

// 指定空值, 必须保证传入的值与默认值类型一致
type emptyValOption struct {
	emptyVal interface{}
}

func (o emptyValOption) apply(opt *whereOptions) {
	opt.emptyVal = o.emptyVal
}

// RequiredOption 标记当前参数必须构建查询, 即使传入的value是空值
func RequiredOption() Option {
	return new(requiredOption)
}

// EmptyValOption用于指定空值, 默认请空空值为 0,""等,
// 但是存在一些特殊请空, 必须空值是-1时, 此时可以通过该方法,指定空值, 当data为该值时,查询条件会被忽略
func EmptyValOption(emptyVal interface{}) Option {
	return emptyValOption{emptyVal: emptyVal}
}

type withSpace string

func (o withSpace) apply(opt *whereOptions) {
	opt.withSpace = string(o)
}

// or 条件
func OrOption() Option {
	return withSpace(OrWithSpace)
}
