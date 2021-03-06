package util

import (
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"toomhub/setting"
)

var DB *gorm.DB
var err error

func MysqlInit() {
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold: time.Second,       // 慢 SQL 阈值
			LogLevel:      gormLogger.Silent, // Log level
			Colorful:      false,             // 禁用彩色打印
		})
	dsn := setting.ZConfig.Database.User + ":" + setting.ZConfig.Database.Password + "@tcp(" + setting.ZConfig.Database.Host + ")/" + setting.ZConfig.Database.Name + "?charset=" + setting.ZConfig.Database.Charset + "&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: newLogger,
	})

	sqlDB, _ := DB.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Second * 100)

	if err != nil {
		fmt.Println(err)
	}
}
