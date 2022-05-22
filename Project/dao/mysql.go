// 定义 Mysql 的数据库配置
// 创建人：吴润泽
// 创建时间：2022-5-15

package dao

import (
	"Project/config"
	"Project/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// InitMysql 初始化 Mysql 数据库连接
func InitMysql() error {
	var db *gorm.DB
	var err error
	// DSN: https://github.com/go-sql-driver/mysql#dsn-data-source-name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.MysqlCfg.UserName,
		config.MysqlCfg.Password,
		config.MysqlCfg.Url,
		config.MysqlCfg.Port,
		config.MysqlCfg.DBName,
		config.MysqlCfg.CharSet,
	)
	db, err = gorm.Open(mysql.New(mysql.Config{
		// DSN Data Source Name
		DSN: dsn,
		// string 类型字段的默认长度
		DefaultStringSize: 256,
		// 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DisableDatetimePrecision: true,
		// 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameIndex: true,
		// 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		DontSupportRenameColumn: true,
		// 根据当前 MySQL 版本自动配置
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 关闭外键约束
		QueryFields:                              true, // select 所有字段而非 select *
	})
	if err != nil {
		log.Fatalf("Connection error: %v\n", err)
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error occured when creating sqlDB: %v\n", err)
		return err
	}
	// 最大闲置连接数
	sqlDB.SetMaxIdleConns(20)
	// 连接池最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 连接池连接最大存活时间，小于 0 表示不关闭
	sqlDB.SetConnMaxLifetime(-1)

	// 自动迁移，在这里添加你的 models
	err = db.AutoMigrate(
		&models.User{},
		&models.Video{},
	)
	if err != nil {
		log.Printf("AutoMigrate error: %v\n", err)
		return err
	}

	DB = db // 暴露在外
	return nil
}
