/*
	ATM is a simple ATM machine cli program for the Alt_School third assignment that has the following features:

	- create an account.
	- change Pin.
	- check account balance
	- withdraw funds
	- deposit funds
	- cancel/exit/logout (selecting this option should exit the program)
*/

package main

import (
	"bufio"
	"errors"
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

type app struct {
	name                       string
	currentCommand             *command
	usage, description, author string
	After                      func(a *app)
	Accounts                   map[string]*account
	Commands                   []*command
}

// account is a user bank account.
type account struct {
	balance        float64
	transactionPin string
}

type command struct {
	action       func(a *app)
	usage        string
	requiredArgs int
}

// newATM returns an new CLI ATM Machine.
func newATM() *app {
	app := &app{
		name:        "CLI ATM Machine",
		usage:       "A CLi ATM Machine with basic bank/ATM features.",
		description: "This is an ATM Machine program that performs basic bank/ATM transactions. Note all amounts are in NGN.",
		author:      "Philemon Ukane",
		After:       performOtherTransaction,
		Accounts: map[string]*account{ // Add default account.
			defaultAccount: {
				balance:        5006267,
				transactionPin: defaultPIN,
			}},
		Commands: []*command{
			{
				action:       createAccount,
				usage:        "To Create a bank account enter username and pin eg. philemon 1234",
				requiredArgs: 2,
			}, {
				action:       depositToAccount,
				usage:        "To deposit an amount into an account, enter the amount, username and pin. e.g 2000 philemon 1234",
				requiredArgs: 3,
			}, {
				action:       withdrawFromAccount,
				usage:        "To withdraw an amount from an account, enter the amount, username and pin. eg 2000 philemon 1234",
				requiredArgs: 3,
			}, {
				action:       checkAccountBalance,
				usage:        "To check an account's balance, enter username and pin. eg philemon 1234",
				requiredArgs: 2,
			}, {
				action:       changeAccountPin,
				usage:        "To change an account's transaction pin enter username, oldpin and newpin, eg philemon 1234 44567",
				requiredArgs: 3,
			}, {
				action:       logOut,
				usage:        "To Shutdown the program",
				requiredArgs: 3,
			}, {
				action: printHelpMsg,
				usage:  "To print this help message.",
			}},
	}

	return app
}

// startAndWelcomeUser prints a welcome message and displays app help message,
// if the default pin is provided.
func startAndWelcomeUser(a *app) {
	var pin string
	for pin != a.Accounts[defaultAccount].transactionPin {
		fmt.Fprintf(os.Stdin, "Enter the pin for the default account, e.g 1234\n")
		pin = readArgs()[0]
	}

	account, err := retrieveAndVerifyAccount(a, defaultAccount, pin)
	if err != nil {
		printErr(a, err)
		return
	}

	fmt.Fprintln(os.Stdin, "--------------------------------------")
	fmt.Fprintf(os.Stdin, "%s, Welcome to the CLI ATM Machine.\n\n", strings.ToUpper(defaultAccount))
	fmt.Fprintf(os.Stdin, "Your current balance is %s.\n", stringAmount(account.balance, naira))
	printHelpMsg(a)
}

// createAccount creates a new bank account for a user.
func createAccount(a *app) {
	args := ensureRequiredArgs(a.currentCommand)

	username := args[0]
	txPin := args[1]

	if username == "" || txPin == "" {
		printErr(a, errors.New("transaction pin or username cannot be empty"))
		return
	}

	_, ok := a.Accounts[username]
	if ok {
		printErr(a, fmt.Errorf("account with username %s already exists", username))
		return
	}

	a.Accounts[username] = &account{
		balance:        2000,
		transactionPin: txPin,
	}

	fmt.Fprintf(os.Stdin, "Account with username %s has been created successfully.\n", username)
	a.After(a)
}

// depositToAccount credits an account with the amount specified via the
// command-line.
func depositToAccount(a *app) {
	args := ensureRequiredArgs(a.currentCommand)

	amt := args[0]
	username := args[1]
	txPin := args[2]

	amount, err := strconv.ParseFloat(amt, 64)
	if err != nil {
		printErr(a, err)
		return
	}

	if amount <= 0 {
		printErr(a, errors.New("deposit amount must be greater than zero"))
		return
	}

	account, err := retrieveAndVerifyAccount(a, username, txPin)
	if err != nil {
		printErr(a, err)
		return
	}
	account.balance += amount

	fmt.Fprintf(os.Stdin, "%s, the deposit to your account was successful.\nYour new balance is %s.\n",
		strings.ToUpper(username), stringAmount(account.balance, naira))
	a.After(a)
}

// withdrawFromAccount debits the amount specified on the command-line from an
// account.
func withdrawFromAccount(a *app) {
	args := ensureRequiredArgs(a.currentCommand)

	amt := args[0]
	username := args[1]
	txPin := args[2]

	amount, err := strconv.ParseFloat(amt, 64)
	if err != nil {
		printErr(a, err)
		return
	}

	if amount <= 0 {
		printErr(a, errors.New("deposit amount must be greater than zero"))
		return
	}

	account, err := retrieveAndVerifyAccount(a, username, txPin)
	if err != nil {
		printErr(a, err)
		return
	}

	if account.balance < amount {
		printErr(a, fmt.Errorf("sorry, your account balance is insufficient for this withdrawal.\nYou have %s and want to withdraw %s",
			stringAmount(account.balance, naira), stringAmount(amount, naira)))
		return
	}
	account.balance -= amount

	fmt.Fprintf(os.Stdin, "%s, Your withdrawal of %s has been concluded successfully.\nYour remaining balance is: %s\n",
		strings.ToUpper(username), stringAmount(amount, naira), stringAmount(account.balance, naira))
	a.After(a)
}

// checkAccountBalance checks and reports an account's balance.
func checkAccountBalance(a *app) {
	args := ensureRequiredArgs(a.currentCommand)

	username := args[0]
	txPin := args[1]
	account, err := retrieveAndVerifyAccount(a, username, txPin)
	if err != nil {
		printErr(a, err)
		return
	}

	fmt.Fprintf(os.Stdin, "%s, Your account balance is %s.\n",
		strings.ToUpper(username), stringAmount(account.balance, naira))
	a.After(a)
}

// changeAccountPin changes an account's transaction pin.
func changeAccountPin(a *app) {
	args := ensureRequiredArgs(a.currentCommand)

	username := args[0]
	oldTxPin := args[1]
	newTxPin := args[3]

	if newTxPin == "" {
		printErr(a, errors.New("new transaction pin cannot be empty"))
		return
	}

	userAccount, err := retrieveAndVerifyAccount(a, username, oldTxPin)
	if err != nil {
		printErr(a, err)
		return
	}

	userAccount.transactionPin = newTxPin

	fmt.Fprintf(os.Stdin, "%s, Your transaction pin has been changed successfully.\n",
		strings.ToUpper(username))
	a.After(a)
}

// logOut exits the CLI ATM program.
func logOut(a *app) {
	fmt.Fprintf(os.Stdin, "Exiting %s .... \n", a.name)
	os.Exit(0)
}

// printHelpMsg prints the help message for an app.
func printHelpMsg(a *app) {
	fmt.Fprintln(os.Stdin, "---------------------------------------")
	fmt.Fprintln(os.Stdin, "About the CLI ATM Machine and Commands")
	fmt.Fprintln(os.Stdin, "---------------------------------------")
	fmt.Println("Command - Usage")
	for key, command := range a.Commands {
		fmt.Fprintf(os.Stdin, "   %d   - %s\n", key, command.usage)
	}
	a.After(a)
}

// retrieveAndVerifyAccount retrieves and verify a user account.
func retrieveAndVerifyAccount(a *app, username, txPin string) (*account, error) {
	if username == "" || txPin == "" {
		return nil, errors.New("transaction pin or username cannot be empty")
	}

	user, ok := a.Accounts[username]
	if !ok {
		return nil, fmt.Errorf("user with username %s does not exist. Please proceed to create an account", username)
	}

	if user.transactionPin != txPin {
		return nil, errors.New("unauthorized: provided transaction pin does not match account's transaction pin")
	}
	return user, nil
}

// performOtherTransaction prompts the user to perform other transactions.
func performOtherTransaction(a *app) {
	fmt.Fprintln(os.Stdin, "---------------------------------------")
	fmt.Println("Do you wish to perform other transactions? For help enter: 6")
	fmt.Println("If you don't wish to continue, exit the app with: 5")
	fmt.Fprintln(os.Stderr, "---------------------------------------")
	var command *command
	for command == nil {
		fmt.Fprintln(os.Stderr, "Enter a valid command")
		args := readArgs()
		for key, cmd := range a.Commands {
			if args[0] == fmt.Sprintf("%d", key) {
				command = cmd
				break
			}
		}
	}

	a.currentCommand = command
	command.action(a)
}

// stringAmount returns an amount with it's currency formatted as a string.
func stringAmount(amount float64, currency string) string {
	return fmt.Sprintf("%.2f %s", amount, currency)
}

// readArgs retrieves user input from the command-line.
func readArgs() []string {
	scanner := bufio.NewScanner(os.Stdin)
	var textInput string
	for scanner.Scan() {
		if textInput = scanner.Text(); textInput != "" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error reading standard input: %v", err))
		return []string{""}
	}

	args := make([]string, 0)
	for _, arg := range strings.SplitAfter(textInput, " ") {
		args = append(args, strings.Trim(arg, " "))
	}

	return args
}

// printErr prints an error and allow for other transactions.
func printErr(a *app, err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	performOtherTransaction(a)
}

// ensureRequiredArgs ensures that the required arguments for a command are
// provided.
func ensureRequiredArgs(cmd *command) []string {
	fmt.Fprintln(os.Stdin, cmd.usage)
	args := readArgs()
	for len(args) != cmd.requiredArgs {
		fmt.Fprintln(os.Stdin, cmd.usage)
		args = readArgs()
	}
	return args
}
