/*
	This is a simple ATM machine cli program for the Alt_School third assignment
	that has the following features:

	- change Pin.
	- check account balance
	- withdraw funds
	- deposit funds
	- cancel/exit/logout (selecting this option should exit the program)
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	// naira is the default and only supported currency.
	naira = "NGN"
	// defaultAccount is a default account that can be used without creating an
	// account.
	defaultAccount = "philemon"
	// defaultPIN is the password fore the default account above.
	defaultPIN = "1234"
)

type atm struct {
	name                string
	description, author string
	user                *account
	Commands            []*command
}

// account is the user account.
type account struct {
	name           string
	balance        float64
	transactionPin string
}

type command struct {
	action func(a *atm)
	usage  string
}

func main() {
	// Create a new atm machine.
	atm := newATM()
	// start the ATM machine
	var pin string
	for pin != atm.user.transactionPin {
		pin = argPrompt("Enter the pin for the default account, i.e 1234")
	}

	fmt.Println("---------------------------------------")
	fmt.Printf("%s, Welcome to the CLI ATM Machine.\n", strings.ToUpper(defaultAccount))
	fmt.Printf("Your current balance is %s.\n", stringAmount(atm.user.balance))
	printHelpMsg(atm)
}

// newATM returns an new CLI ATM Machine.
func newATM() *atm {
	atm := &atm{
		name:        "CLI ATM Machine",
		description: "This is an ATM Machine to performs basic ATM transactions. Note all amounts are in NGN.",
		author:      "Philemon Ukane",
		user: &account{
			name:           defaultAccount,
			balance:        50000,
			transactionPin: defaultPIN,
		},
		Commands: []*command{
			{
				action: depositToAccount,
				usage:  "To deposit an amount.",
			}, {
				action: withdrawFromAccount,
				usage:  "To withdraw an amount.",
			}, {
				action: checkAccountBalance,
				usage:  "To check account balance.",
			}, {
				action: changeAccountPin,
				usage:  "To change transaction pin.",
			}, {
				action: logOut,
				usage:  "To shutdown the machine",
			}, {
				action: printHelpMsg,
				usage:  "To print this help message.",
			}},
	}
	return atm
}

// depositToAccount credits a user with the amount specified via the
// command-line.
func depositToAccount(a *atm) {
	var err error
	var amount float64
	for amount == 0 {
		amt := argPrompt("Enter an amount to deposit")
		amount, err = strconv.ParseFloat(amt, 64)
		if err != nil {
			fmt.Println("error: invalid amount entered")
		}
	}

	a.user.balance += amount

	fmt.Printf("%s, the deposit to your account was successful\n",
		strings.ToUpper(a.user.name))
	performOtherTransactionPrompt(a)
}

// withdrawFromAccount debits the amount specified on the command-line from an
// account.
func withdrawFromAccount(a *atm) {
	var err error
	var amount float64
	for amount == 0 {
		amt := argPrompt("Enter an amount to withdraw")
		amount, err = strconv.ParseFloat(amt, 64)
		if err != nil {
			fmt.Println("error: invalid amount entered")
		}
	}

	if a.user.balance < amount {
		fmt.Println("sorry, your account balance is insufficient for this withdrawal")
		performOtherTransactionPrompt(a)
		return
	}
	a.user.balance -= amount

	fmt.Printf("%s, Your withdrawal of %s is successful.\nYour remaining balance is: %s\n",
		strings.ToUpper(a.user.name), stringAmount(amount), stringAmount(a.user.balance))
	performOtherTransactionPrompt(a)
}

// checkAccountBalance checks and reports an account's balance.
func checkAccountBalance(a *atm) {
	fmt.Fprintf(os.Stdin, "%s, Your account balance is %s.\n",
		strings.ToUpper(a.user.name), stringAmount(a.user.balance))
	performOtherTransactionPrompt(a)
}

// changeAccountPin changes an account's transaction pin.
func changeAccountPin(a *atm) {
	var oldTxPin string
	for oldTxPin == "" || oldTxPin != a.user.transactionPin {
		oldTxPin = argPrompt("Enter your old pin")
		if oldTxPin != a.user.transactionPin {
			fmt.Println("old pin is incorrect")
		}
	}

	var newPin string
	for newPin == "" || newPin == a.user.transactionPin {
		newPin = argPrompt("Enter a new pin")
		if newPin == a.user.transactionPin {
			fmt.Println("cannot use your current pin")
		}
	}

	a.user.transactionPin = newPin

	fmt.Fprintf(os.Stdin, "%s, Your pin has been changed successfully.\n",
		strings.ToUpper(a.user.name))
	performOtherTransactionPrompt(a)
}

// logOut exits the CLI ATM program.
func logOut(a *atm) {
	fmt.Fprintf(os.Stdin, "Exiting %s .... \n", a.name)
	os.Exit(0)
}

// printHelpMsg prints the help message for an app.
func printHelpMsg(a *atm) {
	fmt.Println("---------------------------------------")
	fmt.Println("Available Operations: ")
	fmt.Println("---------------------------------------")
	fmt.Println("Command - Usage")
	for key, command := range a.Commands {
		fmt.Printf("   %d    - %s\n", key, command.usage)
	}
	performOtherTransactionPrompt(a)
}

// performOtherTransactionPrompt prompts the user to perform other transactions.
func performOtherTransactionPrompt(a *atm) {
	fmt.Println("---------------------------------------")
	input := argPrompt("Do you wish to perform other transactions? yes or no")

	switch strings.ToLower(input) {
	case "y", "yes":
		var command *command
		for command == nil {
			input = argPrompt("Choose an operation")
			for key, cmd := range a.Commands {
				if input == fmt.Sprintf("%d", key) {
					command = cmd
					break
				}
			}
		}

		command.action(a)
	case "n", "no":
		logOut(a)
	default:
		performOtherTransactionPrompt(a)
	}
}

// stringAmount returns an amount with formatted as a string.
func stringAmount(amount float64) string {
	return fmt.Sprintf("%.2f %s", amount, naira)
}

// argPrompt retrieves user input from the command-line.
func argPrompt(message string) string {
	scanner := bufio.NewScanner(os.Stdin)
	var textInput string
	for textInput == "" {
		fmt.Println(message)
		scanner.Scan()
		textInput = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading standard input: %v\n", err)
		return ""
	}
	return textInput
}
