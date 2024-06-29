package conv

import (
	"errors"
	"fmt"
	"reflect"
)

// Mapper 用于对象隐射, 将一个对象从一个类型转换为另一个类型
// 支持单个对象或者切片, 同时如果字段为结构体, 会遍历子结构体匹配
//
// 默认情况, 通过源对象与目标对象的名称匹配(也可通过方法 Mapper.Bind 来指定字段关联)
// 匹配时, 支持部分数据类型相互转换
// 当前支持的类型转换:
//
//		Int Int8 Int16 Int32 Int64 Uint Uint8 Uint16 Uint32 Uint64 --> Float32
//	 Int Int8 Int16 Int32 Int64 Uint Uint8 Uint16 Uint32 Uint64 --> Float64
//	 Int64 --> Int32
//	 Int32 --> Int64
type Mapper struct {
	bind     map[string]string
	excludes map[string]struct{}
}

func NewMapper() *Mapper {
	return &Mapper{}
}

// Bind 绑定隐射字段,默认情况, 使用名称一致来做隐射, 使用者可以手动指定隐射绑定的字段,
// key为目标对象的字段名称, value为源的字段名称
func (m *Mapper) Bind(bind map[string]string) *Mapper {
	m.bind = bind
	return m
}

// Exclude 指定目标对象需要排除的字段, 指定后, 这些字段回本忽略
func (m *Mapper) Exclude(excludes ...string) *Mapper {
	m.excludes = make(map[string]struct{})

	for _, v := range excludes {
		m.excludes[v] = struct{}{}
	}
	return m
}

// Map 隐射对象
// @param source 数据源对象
// @param dest 目标对象,必须传入指针
func (m *Mapper) Map(source, dest interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	var destType = reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errors.New("dest 参数必须为指针类型！")
	}
	var sourceVal = reflect.ValueOf(source)
	var destVal = reflect.ValueOf(dest).Elem()

	m.mapValues(sourceVal, destVal, true)
	return nil
}

func (m *Mapper) mapValues(sourceVal, destVal reflect.Value, loose bool) {
	destType := destVal.Type()
	if destType.Kind() == reflect.Struct {
		if sourceVal.Type().Kind() == reflect.Ptr {
			if sourceVal.IsNil() {
				// If source is nil, it maps to an empty struct
				sourceVal = reflect.New(sourceVal.Type().Elem())
			}
			sourceVal = sourceVal.Elem()
		}
		for i := 0; i < destVal.NumField(); i++ {
			m.mapField(sourceVal, destVal, i, loose)
		}
	} else if m.isTypeValidate(sourceVal.Type(), destType) {
		m.setValue(sourceVal, destVal)
	} else if destType.Kind() == reflect.Ptr {
		if m.valueIsNil(sourceVal) {
			return
		}
		val := reflect.New(destType.Elem())
		m.mapValues(sourceVal, val.Elem(), loose)
		destVal.Set(val)
	} else if destType.Kind() == reflect.Slice {
		// 判断类型要一致
		if sourceVal.Type().Kind() != reflect.Slice {
			panic("目标类型是切片,源数据类型也必须是切片类型")
		}
		m.mapSlice(sourceVal, destVal, loose)
	} else {
		panic(fmt.Sprintf("暂不支持类型(%v)转换为类型(%v)", sourceVal.Type().Kind(), destType.Kind()))
	}
}

// 源类型隐射,这些类型允许转换为对应的其他类型
func (m *Mapper) isTypeValidate(sourceType, destType reflect.Type) bool {
	if sourceType == destType {
		return true
	}
	// 针对目标是float64做特殊处理(因为我们很多业务数据库存储的是int32|int64,但是实际的业务需要的是float64)
	switch destType.Kind() {
	case reflect.Float32:
		if sourceType.Kind() >= reflect.Int && sourceType.Kind() <= reflect.Uint64 {
			return true
		}
	case reflect.Float64:
		// 针对数字类型(float32)的,都可以转换为float64
		if sourceType.Kind() >= reflect.Int && sourceType.Kind() <= reflect.Float32 {
			return true
		}
		// 允许int32和int64相互转换
	case reflect.Int32:
		if sourceType.Kind() == reflect.Int || sourceType.Kind() == reflect.Int64 {
			return true
		}
		// 允许int32和int64相互转换
	case reflect.Int64:
		if sourceType.Kind() == reflect.Int || sourceType.Kind() == reflect.Int32 {
			return true
		}
	case reflect.String:
		// 允许将 *string转换为string
		if sourceType.Kind() == reflect.Ptr && sourceType.Elem().Kind() == reflect.String {
			return true
		}
	}

	return false
}

func (m *Mapper) setValue(sourceVal, destVal reflect.Value) {
	sourceType := sourceVal.Type()
	destType := destVal.Type()
	// 如果类型一样,直接赋值
	if sourceType == destType {
		destVal.Set(sourceVal)
		return
	}

	// 类型不一致时, 做兼容操作
	switch destType.Kind() {
	case reflect.Float32:
		if sourceType.Kind() >= reflect.Int && sourceType.Kind() <= reflect.Int64 {
			destVal.Set(reflect.ValueOf(float32(sourceVal.Int())))
			return
		}
		if sourceType.Kind() >= reflect.Uint && sourceType.Kind() <= reflect.Uint64 {
			destVal.Set(reflect.ValueOf(float32(sourceVal.Uint())))
			return
		}
	// 针对目标是float64做特殊处理(因为我们很多业务数据库存储的是int32|int64,但是实际的业务需要的是float64)
	case reflect.Float64:
		if sourceType.Kind() >= reflect.Int && sourceType.Kind() <= reflect.Int64 {
			destVal.Set(reflect.ValueOf(float64(sourceVal.Int())))
			return
		}
		if sourceType.Kind() >= reflect.Uint && sourceType.Kind() <= reflect.Uint64 {
			destVal.Set(reflect.ValueOf(float64(sourceVal.Uint())))
			return
		}
		if sourceType.Kind() >= reflect.Float32 && sourceType.Kind() <= reflect.Float64 {
			destVal.Set(reflect.ValueOf(sourceVal.Float()))
			return
		}
	// int64和int32相互转换
	case reflect.Int64:
		destVal.Set(reflect.ValueOf(sourceVal.Int()))
	case reflect.Int32:
		destVal.Set(reflect.ValueOf(int32(sourceVal.Int())))
	case reflect.String:
		if !sourceVal.IsNil() {
			destVal.Set(reflect.ValueOf(sourceVal.Elem().String()))
		}
	}
}

func (m *Mapper) mapSlice(sourceVal, destVal reflect.Value, loose bool) {
	destType := destVal.Type()
	length := sourceVal.Len()
	target := reflect.MakeSlice(destType, length, length)
	for j := 0; j < length; j++ {
		val := reflect.New(destType.Elem()).Elem()
		m.mapValues(sourceVal.Index(j), val, loose)
		target.Index(j).Set(val)
	}

	if length == 0 {
		m.verifyArrayTypesAreCompatible(sourceVal, destVal, loose)
	}
	destVal.Set(target)
}

func (m *Mapper) verifyArrayTypesAreCompatible(sourceVal, destVal reflect.Value, loose bool) {
	dummyDest := reflect.New(reflect.PtrTo(destVal.Type()))
	dummySource := reflect.MakeSlice(sourceVal.Type(), 1, 1)
	m.mapValues(dummySource, dummyDest.Elem(), loose)
}

func (m *Mapper) mapField(source, destVal reflect.Value, i int, loose bool) {
	destType := destVal.Type()
	destFieldName := destType.Field(i).Name
	// 如果是排除字段, 则直接返回
	if _, ok := m.excludes[destFieldName]; ok {
		return
	}
	sourceFieldName := destFieldName
	if m.bind != nil {
		if v, ok := m.bind[destFieldName]; ok {
			sourceFieldName = v
		}
	}

	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("mapping错误,隐射错误字段:%s 目标对象:%v, 源对象:%v Error: %v", destFieldName, destType, source.Type(), r))
		}
	}()

	destField := destVal.Field(i)
	if destType.Field(i).Anonymous {
		m.mapValues(source, destField, loose)
	} else {
		if m.valueIsContainedInNilEmbeddedType(source, sourceFieldName) {
			return
		}
		sourceField := source.FieldByName(sourceFieldName)
		if (sourceField == reflect.Value{}) {
			if loose {
				return
			}
			if destField.Kind() == reflect.Struct {
				m.mapValues(source, destField, loose)
				return
			} else {
				for i := 0; i < source.NumField(); i++ {
					if source.Field(i).Kind() != reflect.Struct {
						continue
					}
					if sourceField = source.Field(i).FieldByName(sourceFieldName); (sourceField != reflect.Value{}) {
						break
					}
				}
			}
		}
		m.mapValues(sourceField, destField, loose)
	}
}

func (m *Mapper) valueIsNil(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Ptr && value.IsNil()
}

func (m *Mapper) valueIsContainedInNilEmbeddedType(source reflect.Value, fieldName string) bool {
	structField, _ := source.Type().FieldByName(fieldName)
	ix := structField.Index
	if len(structField.Index) > 1 {
		parentField := source.FieldByIndex(ix[:len(ix)-1])
		if m.valueIsNil(parentField) {
			return true
		}
	}
	return false
}
