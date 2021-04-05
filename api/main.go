package main

import (
	"github.com/davidwarshaw/golang-user-crud/api/server"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("port", "8080")
	server.Setup().Run(":" + viper.GetString("port"))
}
