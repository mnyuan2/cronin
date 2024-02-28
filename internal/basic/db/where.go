package db

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	OpEqual     = "="
	OpIn        = "IN"
	OpLike      = "LIKE"
	OpFindInSet = "FIND_IN_SET"
	OpGt        = ">"
	OpGte       = ">="
	OpLt        = "<"
	OpLte       = "<="
)

// Where 查询条件构建器
type Where struct {
	wheres []string
	values []interface{}

	// 其它条款
	clauses map[string]interface{}
}

func NewWhere() *Where {
	return &Where{}
}

// 原始的查询
// @query 原始的查询条件, 必须带?号,开头不用 AND, 支持多个查询条件组合, 例如 name = ? AND id > ?
// @args 参数的值
func (builder *Where) Raw(query string, args ...interface{}) *Where {
	builder.wheres = append(builder.wheres, query)
	if len(args) > 0 {
		builder.values = append(builder.values, args...)
	}

	return builder
}

// FindInSet 用于查询数据库字段的值是逗号隔开的数据, 当传入的参数value是切片或数据(或逗号隔开字符串), 只要匹配一个值, 就能取到数据
// 注意, 该查询条件性能很低, 数据量较大时谨慎使用
// @param field 字段
// @param value 数据值, value可传入基础类型的切片或者数组,或者逗号隔开的字符串,当值为空值时,忽略该查询条件
func (builder *Where) FindInSet(field string, value interface{}, options ...Option) *Where {
	opt := ApplyOptions(options...)
	kind := reflect.TypeOf(value).Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		valueOf := reflect.ValueOf(value)
		length := valueOf.Len()
		if opt.required || length > 0 {
			where := " ("
			for i := 0; i < length; i++ {
				item := valueOf.Index(i)
				where += fmt.Sprintf("FIND_IN_SET(%v,%v)", item, field)
				// 表示不是最后一个, 添加 or
				if i < length-1 {
					where += " OR "
				}
			}
			where += " )"

			builder.wheres = append(builder.wheres, where)
		}
	} else {
		if opt.required || !opt.isZero(value) {
			if val, ok := value.(string); ok {
				items := strings.Split(val, ",")
				length := len(items)
				if length > 0 {
					where := " ("
					for i := 0; i < length; i++ {
						item := items[i]
						where += fmt.Sprintf("FIND_IN_SET(%v,%v)", item, field)
						// 表示不是最后一个, 添加 or
						if i < length-1 {
							where += " OR "
						}
					}
					where += " )"

					builder.wheres = append(builder.wheres, where)
				}
			} else {
				// 单个值的情况
				builder.wheres = append(builder.wheres, fmt.Sprintf(" AND FIND_IN_SET(%v,%v)", value, field))
			}
		}
	}

	return builder
}

// 构建json 路径in查询
func (builder *Where) JsonPathIn(field string, values interface{}, options ...Option) *Where {
	//opt := ApplyOptions(options...)

	args := []string{}
	switch values.(type) {
	case []int32:
		val := values.([]int32)
		for _, v := range val {
			args = append(args, strconv.Itoa(int(v)))
		}
	case []string:
		args = values.([]string)
	default:
		panic("未支持的数据类型")
	}

	if len(args) > 0 {
		str := ""
		for _, v := range args {
			str += "'$.\"" + v + "\"',"
		}

		//_where := fmt.Sprintf("JSON_CONTAINS_PATH(%s, 'one', %s)", field, strings.Trim(str, ","))
		//builder.wheres = append(builder.wheres, _where)
	}
	return builder
}

// JsonIndexEq json查询kv相等 查询
//
//	示例： json_search(tags_key, 'one', 'config_id') = json_search(tags_val, 'one', '8')
func (builder *Where) JsonIndexEq(keyField, valField string, key, val any, options ...Option) *Where {
	builder.Raw(fmt.Sprintf("json_search(%s, 'one', ?) = json_search(%s,'one', ?)", keyField, valField), key, val)
	return builder
}

// 当传入的Value为空值时,忽略该查询条件
func (builder *Where) Equal(field string, value interface{}, options ...Option) *Where {
	return builder.op(field, OpEqual, value, options...)
}

// Equal·等于 的简写方法
func (builder *Where) Eq(field string, value interface{}, options ...Option) *Where {
	return builder.op(field, OpEqual, value, options...)
}

// 当传入的Value为空值时,忽略该查询条件
func (builder *Where) Like(field string, value interface{}, options ...Option) *Where {
	return builder.op(field, OpLike, value, options...)
}

// value可传入基础类型的切片或者数组,或者逗号隔开的字符串,当值为空值时,忽略该查询条件
func (builder *Where) In(field string, value interface{}, options ...Option) *Where {
	kind := reflect.TypeOf(value).Kind()
	// 切片和数组
	if kind == reflect.Slice || kind == reflect.Array {
		return builder.op(field, OpIn, value, options...)
	} else if kind == reflect.String {
		return builder.op(field, OpIn, strings.Split(value.(string), ","), options...)
	} else {
		panic("不支持的数据类型,In只允许传入数组或者切片,或者逗号隔开的字符串")
	}
}

// 当传入的Value为空值时,忽略该查询条件
func (builder *Where) Gt(field string, value interface{}, options ...Option) *Where {
	return builder.op(field, OpGt, value, options...)
}

// 当传入的Value为空值时,忽略该查询条件
func (builder *Where) Gte(field string, value interface{}, options ...Option) *Where {
	return builder.op(field, OpGte, value, options...)
}

// 当传入的Value为空值时,忽略该查询条件
func (builder *Where) Lt(field string, value interface{}, options ...Option) *Where {
	return builder.op(field, OpLt, value, options...)
}

// 当传入的Value为空值时,忽略该查询条件
func (builder *Where) Lte(field string, value interface{}, options ...Option) *Where {
	return builder.op(field, OpLte, value, options...)
}

// 开始构建, 生成查询语句以及参数
func (builder *Where) Build() (whereStr string, args []interface{}) {
	// 固定加上1=1， 防止外部查询条件还需要判断
	builder.wheres = append([]string{"1 = 1"}, builder.wheres...)
	return strings.Join(builder.wheres, " AND "), builder.values
}

// 将数据添加到 where 条件
// @param field 字段
// @param op 操作, 当前仅支持 = IN (not in) 等简单的类型
// @param value 数据值, 仅支持基础类型,以及简单类型的切片(不做类型判断)
// @param required 是否必须的查询条件，如果为true, 将不会判断value空值， 必须会带上该查询条件;
func (builder *Where) op(field string, op string, value interface{}, options ...Option) *Where {
	if op == OpFindInSet {
		return builder.FindInSet(field, value, options...)
	}

	opt := ApplyOptions(options...)
	kind := reflect.TypeOf(value).Kind()
	// 切片和数组
	if kind == reflect.Slice || kind == reflect.Array {
		valueOf := reflect.ValueOf(value)
		if opt.required || valueOf.Len() > 0 {
			builder.wheres = append(builder.wheres, fmt.Sprintf("%v %v ?", field, op))
			builder.values = append(builder.values, value)
		}
	} else {
		if opt.required || !opt.isZero(value) {
			appendVal := value
			if op == OpLike {
				if value, ok := value.(string); ok {
					appendVal = "%" + value + "%"
				}
			}

			builder.wheres = append(builder.wheres, fmt.Sprintf("%v %v ?", field, op))
			builder.values = append(builder.values, appendVal)
		}
	}

	return builder
}

// 返回数
func (builder *Where) Limit(limit int) *Where {
	builder.clauses["LIMIT"] = limit
	return builder
}

// 数据偏移位
func (builder *Where) Offset(limit int) *Where {
	builder.clauses["OFFSET"] = limit
	return builder
}
