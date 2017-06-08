package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/conradludgate/ladgate/module"
	flag "github.com/spf13/pflag"
)

func main() {
	protocol := flag.String("protocol", "", "Ladgate Protocol to listen on")

	url, name, pass := module.LoadConfig("admin")
	m := module.NewModule(url, name, pass)

	if *protocol == "" {
		fmt.Println("Please provide a protocol")
		return
	}

	go func() {
		log.Fatal(m.Listen(*protocol, nil))
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for fmt.Print(name + ": "); scanner.Scan(); fmt.Print(name + ": ") {
		m.SendMessage(*protocol, scanner.Text())
	}
}
