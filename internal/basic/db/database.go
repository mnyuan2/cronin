package db

import (
	"context"
	"cron/internal/basic/config"
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

var (
	_db  *gorm.DB
	once sync.Once
)

// 连接数据库
func New(ctx context.Context) *MyDB {
	once.Do(func() {
		conf := config.DbConf()
		switch conf.Driver {
		case "mysql":
			if _db = Conn(conf.Mysql); _db.Error != nil {
				panic(_db.Error)
			}
		default:
			panic(fmt.Sprintf("database.driver=%s 为支持", conf.Driver))
		}

	})

	// 根据实例,修改上下文
	return &MyDB{_db.WithContext(ctx)}

}

func Conn(conf *config.MysqlSource) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=false&loc=Local",
		conf.Username, conf.Password, conf.Hostname, conf.Port, conf.Database)
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "", // 表前缀
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		db.AddError(err)
	} else if err = polling(db); err != nil { // 启用连接池
		db.AddError(fmt.Errorf("连接池设置异常 %w", err))
	} else if conf.Debug { // 调试模式
		db = db.Debug()
	}

	return db
}

// 连接 clickhouse
func ConnClickhouse(conf *config.MysqlSource) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=false&loc=Local",
		conf.Username, conf.Password, conf.Hostname, conf.Port, conf.Database)
	// 连接数据库
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		db.AddError(err)
	} else if err = polling(db); err != nil { // 启用连接池
		db.AddError(fmt.Errorf("连接池设置异常 %w", err))
	} else if conf.Debug { // 调试模式
		db = db.Debug()
	}

	return db
}

// 设置程序池;
// 这个有空要研究一下
func polling(_db *gorm.DB) error {
	sqlDb, err := _db.DB()
	if err != nil {
		return err
	}
	// TODO: 这些参数要迁移到配置文件中
	// 设置连接池;
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDb.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDb.SetConnMaxLifetime(time.Hour)
	return nil
}
