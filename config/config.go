package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	MysqlUser     string
	MysqlPassword string
	MysqlPort     string
	MysqlDb       string
	MysqlHost     string
	RedisHost     string
	RedisPort     string
)

func init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err.Error())

	}
	MysqlHost = viper.Get("mysql.host").(string)
	MysqlUser = viper.Get("mysql.user").(string)
	MysqlPassword = viper.Get("mysql.password").(string)
	MysqlPort = viper.Get("mysql.port").(string)
	MysqlDb = viper.Get("mysql.db").(string)
	RedisHost = viper.Get("redis.host").(string)
	RedisPort = viper.Get("redis.port").(string)
	fmt.Println(MysqlPort, MysqlUser, MysqlDb, MysqlPassword, MysqlHost, RedisPort, RedisHost)
}
