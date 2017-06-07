package main

import (
	"strconv"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadConfig() (addr, cert, key string) {
	config := flag.StringP("config", "c", "ladgate", "Config file to find connection info")

	flag.IntP("port", "p", 2367, `Port to listen on`)

	flag.String("cert", "/etc/ladgate/ssl/cert.pem", "SSL Cert")
	flag.String("key", "/etc/ladgate/ssl/key.pem", "SSL Key")

	flag.Parse()

	viper.SetConfigName(*config)

	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("$HOME/.config/ladgate")
	viper.AddConfigPath("$HOME/.ladgate")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	viper.BindPFlag("port", flag.Lookup("port"))
	viper.BindPFlag("cert", flag.Lookup("cert"))
	viper.BindPFlag("key", flag.Lookup("key"))

	viper.ReadInConfig()

	return ":" + strconv.Itoa(viper.GetInt("port")), viper.GetString("cert"), viper.GetString("key")
}
