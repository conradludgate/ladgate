package main

import (
	"fmt"
	"log"

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

	module.HandleFunc(*protocol, HandleMessage)

	log.Fatal(m.Listen(*protocol, nil))
}

func HandleMessage(m *module.Module, msg module.Message) {
	fmt.Printf("%s | %s\t | %s", msg.Time.Format("2006-01-02 15:04:05"), msg.Module, msg.Data)
}
