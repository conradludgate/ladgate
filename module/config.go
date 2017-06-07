package module

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/bgentry/speakeasy"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadConfig(name string) (URL string, Name string, password string) {
	config := flag.StringP("config", "c", "module", "Config file to find connection info")

	moduleName := flag.StringP("name", "n", name, "Username of the module")

	flag.StringP("url", "u", "", `URL to connect to (default "127.0.0.1")`)
	flag.IntP("port", "p", 0, `Port to connect to (default 2367)`)

	flag.String("pass", "", "Password for the module")

	flag.Parse()

	viper.SetConfigName(*config)

	viper.AddConfigPath("/etc/ladgate")
	viper.AddConfigPath("/etc/ladgate/" + *moduleName)
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("$HOME/.config/ladgate")
	viper.AddConfigPath("$HOME/.config/ladgate/" + *moduleName)
	viper.AddConfigPath("$HOME/.ladgate")
	viper.AddConfigPath("$HOME/.ladgate/" + *moduleName)
	viper.AddConfigPath(".")
	viper.AddConfigPath(*moduleName)
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/" + *moduleName)

	//viper.RegisterAlias("url", *moduleName+".url")
	//viper.RegisterAlias("port", *moduleName+".port")

	viper.BindPFlag(*moduleName+".url", flag.Lookup("url"))
	viper.BindPFlag(*moduleName+".port", flag.Lookup("port"))
	viper.BindPFlag(*moduleName+".password", flag.Lookup("pass"))

	if viper.ReadInConfig() == viper.ConfigFileNotFoundError {
		viper.SetConfigName(*moduleName)
		viper.ReadInConfig()
	}

	URL = viper.GetString(*moduleName + ".url")
	if URL == "" {
		URL = viper.GetString("url")
	}
	if URL == "" {
		URL = "127.0.0.1"
	}

	port := viper.GetInt(*moduleName + ".port")
	if port == 0 {
		port = viper.GetInt("port")
	}
	if port == 0 {
		port = 2367
	}

	u, err := url.Parse(URL)
	if err != nil {
		panic(fmt.Errorf("Fatal error url: %s \n", err))
	}

	if u.Scheme == "" {
		u, _ = url.Parse("ws://" + URL)
	}

	if u.Port() == "" {
		u.Host += ":" + strconv.Itoa(port)
	}

	if u.Path == "" {
		u.Path = "ws"
	}

	password = viper.GetString(*moduleName + ".password")
	if password == "" {
		password = viper.GetString(*moduleName + ".pass")
	}
	if password == "" {
		password = viper.GetString("password")
	}
	if password == "" {
		password = viper.GetString("pass")
	}
	if password == "" {
		password, err = speakeasy.Ask("Password for " + *moduleName + ": ")
		if err != nil {
			panic(err)
		}
	}

	return u.String(), *moduleName, password
}
