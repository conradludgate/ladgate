package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/conradludgate/ladgate/module"
	flag "github.com/spf13/pflag"
)

func main() {
	protocol := flag.String("protocol", "", "Ladgate Protocol to listen on")

	url, name, pass := module.LoadConfig("admin")

	if *protocol == "" {
		fmt.Println("Please provide a protocol")
		return
	}

	m := module.NewModule(url, name, pass)

	err := m.Connect(*protocol, nil)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for fmt.Print(name + ": "); scanner.Scan(); fmt.Print(name + ": ") {
		io.WriteString(m, scanner.Text())
	}
}
