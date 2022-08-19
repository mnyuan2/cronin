package db

import "gorm.io/gorm"

type MyDB struct {
	*gorm.DB
}

// 单表分页查询(由于需要排序,Page方法弃用)
// @params
// @out: 分页数据，需要传入指针
// @page: 当前页，从1开始
// @pageSize: 页码大小，默认10
// @tableName： 表名称
// @fields: 要查的字段,如果为空表示查全部(*)
// @where: 查询条件
// @values: 条件值
// @returns
// @count: 记录数
// @err: 错误
func (gdb *MyDB) Paginate(out interface{}, page int, pageSize int, tableName string, fields string, orderBy string, where string, values ...interface{}) (count int64, err error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if fields == "" {
		fields = "*"
	}

	// 获取数量;
	if err := gdb.Table(tableName).
		Where(where, values...).
		Count(&count).Error; err != nil {
		return 0, err
	}
	offset := (page - 1) * pageSize
	if int64(offset) > count {
		return
	}

	txDb := gdb.Table(tableName).
		Select(fields).
		Where(where, values...)
	if orderBy != "" {
		txDb = txDb.Order(orderBy)
	}

	if err = txDb.Offset(int(offset)).Limit(int(pageSize)).Find(out).Error; err != nil {
		return 0, err
	}
	return
}
