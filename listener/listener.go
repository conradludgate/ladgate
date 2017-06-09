package main

import (
	"fmt"

	"github.com/conradludgate/ladgate/module"
	flag "github.com/spf13/pflag"
)

func main() {
	protocol := flag.String("protocol", "", "Ladgate Protocol to listen on")
	m := module.NewModule(module.LoadConfig("admin"))

	if *protocol == "" {
		fmt.Println("Please provide a protocol")
		return
	}

	module.HandleFunc(module.MatchAllPattern, HandleMessage)

	err := m.Connect(*protocol, nil)
	if err != nil {
		panic(err)
	}

	for {
	}
}

func HandleMessage(m *module.Module, msg module.Message) {
	fmt.Printf("%s | %s\t | %s\n", msg.Time.Format("15:04:05"), msg.Module, msg.Data)
}
