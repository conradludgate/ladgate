package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	irc "github.com/fluffle/goirc/client"
	"github.com/spf13/viper"
)

func main() {
	Connections = make(map[string]*irc.Conn)
	LoadCoreSettings()
	<-LoadCore()
}

func LoadCoreSettings() {
	viper.SetDefault("Server", "localhost:6667")
	viper.SetDefault("SSL", false)

	viper.SetDefault("Nick", "ladgate")
	viper.SetDefault("Ident", "Ladgate IRC Interface")
	viper.SetDefault("Name", "https://github.com/conradludgate/ladgate")

	viper.SetConfigName("config")

	viper.AddConfigPath("/etc/ladgate")
	viper.AddConfigPath("$HOME/.config/ladgate")
	viper.AddConfigPath("$HOME/.ladgate")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}
}

func LoadCore() chan bool {
	cfg := irc.NewConfig("core", "IRC Bot core",
		"https://github.com/conradludgate/ladgate")

	cfg.Server = viper.GetString("Server")
	cfg.Pass = viper.GetString("Pass")

	cfg.SSL = viper.GetBool("SSL")

	if cfg.SSL {
		cfg.SSLConfig = &tls.Config{ServerName: strings.Split(cfg.Server, ":")[0]}
	}

	core := irc.Client(cfg)
	core.EnableStateTracking()

	quit := make(chan bool)
	core.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		time.Sleep(time.Second * 10)
		core.Connect()
	})

	core.HandleFunc(irc.INVITE, CoreOnInvite)
	core.HandleFunc(irc.PRIVMSG, CoreOnPrivMsg)

	if err := core.Connect(); err != nil {
		fmt.Println(err.Error())
		quit <- true
	}

	return quit
}
