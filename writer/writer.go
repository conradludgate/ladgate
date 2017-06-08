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

	go log.Fatal(m.Listen(*protocol, nil))

	for {
		fmt.Print(flag.Lookup("name").Value.String() + ": ")
		var msg string
		fmt.Scanln(&msg)

		m.SendMessage(*protocol, msg)
	}
}
