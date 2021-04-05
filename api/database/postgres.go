package database

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
)

var PK_ERROR_CODE = "ERROR #23505"

func Middleware() gin.HandlerFunc {
	// We need tcp to go across containers
	viper.SetDefault("db_network", "tcp")
	// docker compose DB host
	viper.SetDefault("db_addr", "db:5432")
	// We'll use the superuser default db to simplify seeding
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_password", "postgres")
	viper.SetDefault("db_database", "postgres")

	options := pg.Options{
		Network:  viper.GetString("db_network"),
		Addr:     viper.GetString("db_addr"),
		User:     viper.GetString("db_user"),
		Password: viper.GetString("db_password"),
		Database: viper.GetString("db_database"),
	}

	db := pg.Connect(&options)

	return func(c *gin.Context) {
		c.Set("DB", db)
	}
}
