package main

import (
	"fmt"
	"os"

	"github.com/ukane-philemon/cli-atm/atm"
)

func main() {
	fmt.Println("Enter the pin for the default account in this format, i.e: start --pin 1234")
	args := os.Args
	if len(args) < 2 {
		args = atm.ReadArgs()
	}

	err := atm.StartATM(args)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
