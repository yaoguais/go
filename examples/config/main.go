package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "", "config file")
	flag.Parse()

	viper.SetConfigFile(configFile)
	viper.ReadInConfig()
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if viper.GetString("mysql.host") == "" {
		panic("config error")
	}
}

func main() {
	fmt.Printf("previous mysql.host = %s\n", viper.GetString("mysql.host"))
	os.Setenv("MYSQL_HOST", "192.168.1.1")
	fmt.Printf("after mysql.host = %s\n", viper.GetString("mysql.host"))

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %v\n", e)
		fmt.Printf("reload mysql.host = %s\n", viper.GetString("mysql.host"))
	})

	go func() {
		for {
			time.Sleep(10 * time.Second)
			fmt.Printf("timer load mysql.host = %s\n", viper.GetString("mysql.host"))
			fmt.Printf("os load mysql.host = %s\n", os.Getenv("MYSQL_HOST"))
		}
	}()

	// shell command "export MYSQL_HOST=127.0.0.3" doesn't work

	time.Sleep(86400 * time.Second)
}
