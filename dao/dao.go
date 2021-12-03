package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"software_bottle/config"
)

var (
	DB         *gorm.DB
	RedisCache *Cache
)

func InitRedis() {
	RedisCache = NewRedisCache(0, fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort), FOREVER)
}
func InitMySQL() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MysqlUser, config.MysqlPassword, config.MysqlHost, config.MysqlPort, config.MysqlDb)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}
