package main

import (
	"fmt"
	"os"

	"github.com/ukane-philemon/cli-atm/atm"
)

func main() {
	args := os.Args
	if len(args) < 3 { // If the user provides some arguments, use it. Else prompt to start the ATM Machine.
		fmt.Println("Enter the pin for the default account in this format, i.e: start --pin 1234")
		args = atm.ReadArgs()
	}

	err := atm.StartATM(args)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
