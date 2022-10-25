/*
	ATM is a simple ATM machine cli program for the Alt_School third assignment that has the following features:

	- create an account.
	- change Pin.
	- check account balance
	- withdraw funds
	- deposit funds
	- cancel/exit/logout (selecting this option should exit the program)
*/

package atm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	// naira is the default and only supported currency.
	naira = "NGN"
	// amountFlag is on of the required command-line argument to process a
	// deposit or withdrawal transaction.
	amountFlag = []cli.Flag{&cli.Float64Flag{Name: "amount", Required: true, Usage: "Amount to deposit or withdraw."}}
	// requiredFlags these are the command-line arguments to required to carry
	// out a transaction on an account.
	requiredFlags = []cli.Flag{
		&cli.StringFlag{Name: "username", Value: "test", Required: true, Usage: "The username of the account."},
		&cli.StringFlag{Name: "pin", Value: "1234", Required: true, Usage: "The transaction pin of the account."},
	}
)

// account is a user bank account.
type account struct {
	balance        float64
	transactionPin string
}

// StartATM starts the CLI ATM program.
func StartATM(args []string) error {
	app := cli.NewApp()
	app.Name = "CLI ATM Machine"
	app.Usage = "A CLi ATM Machine with basic bank/ATM features."
	app.Description = "This is an ATM Machine program that performs basic bank/ATM transactions. Note all amounts are in NGN."
	app.Authors = []*cli.Author{{Name: "Philemon Ukane"}}
	// Add a  default user account.
	app.Metadata = map[string]interface{}{
		"philemon": &account{
			transactionPin: "1234",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:      "start",
			Aliases:   []string{"s"},
			Usage:     "To launch the CLI ATM Machine",
			UsageText: "e.g start --pin=1234",
			After:     performOtherTransaction,
			Action:    welcomeUser,
			Flags:     []cli.Flag{&cli.StringFlag{Name: "pin", Value: "1234", Required: true, Usage: "The transaction pin of the default account."}},
		},
		{
			Name:      "createaccount",
			Aliases:   []string{"ca"},
			Usage:     "Create a bank account. Note: if an account already exists, no further action will be carried out.",
			UsageText: "e.g createaccount --username=philemon --pin=1234",
			After:     performOtherTransaction,
			Action:    createAccount,
			Flags:     requiredFlags,
		}, {
			Name:      "deposit",
			Aliases:   []string{"d"},
			Usage:     "To deposit an amount into an account.",
			UsageText: "Specify the amount followed by the username and transaction pin.e.g deposit --amount=2000 --username=philemon --pin=1234",
			After:     performOtherTransaction,
			Action:    depositToAccount,
			Flags:     append(amountFlag, requiredFlags...),
		}, {
			Name:      "withdraw",
			Aliases:   []string{"w"},
			Usage:     "To withdraw an amount from an account.",
			UsageText: "Specify the withdraw amount followed by the username and transaction pin.e.g withdraw --amount=2000 --username=philemon --pin=1234",
			After:     performOtherTransaction,
			Action:    withdrawFromAccount,
			Flags:     append(amountFlag, requiredFlags...),
		}, {
			Name:      "balance",
			Aliases:   []string{"b"},
			Usage:     "To check an account's balance",
			UsageText: "e.g balance --username=philemon --pin=1234",
			After:     performOtherTransaction,
			Action:    checkAccountBalance,
			Flags:     requiredFlags,
		}, {
			Name:      "changepin",
			Aliases:   []string{"cp"},
			Usage:     "To change an account's transaction pin.",
			UsageText: "e.g changepin --username=philemon --pin=1234 --newpin=3456",
			After:     performOtherTransaction,
			Action:    changeAccountPin,
			Flags:     append([]cli.Flag{&cli.StringFlag{Name: "newpin", Required: true, Usage: "New account transaction pin."}}, requiredFlags...),
		}, {
			Name:      "logout",
			Aliases:   []string{"exit", "cancel"},
			Usage:     "To Shutdown the program",
			UsageText: "e.g logout",
			Action:    logOut,
		}}

	err := app.Run(args)
	if err != nil {
		return err
	}

	return nil
}

// welcomeUser prints a welcome message and displays app help message.
func welcomeUser(c *cli.Context) error {
	fmt.Fprintln(c.App.Writer, "Welcome to the CLI ATM Machine")
	return cli.ShowAppHelp(c)
}

// createAccount creates a new bank account for a user.
func createAccount(c *cli.Context) error {
	username := c.String("username")
	txPin := c.String("pin")

	if username == "" || txPin == "" {
		fmt.Fprintln(os.Stderr, "transaction pin or username cannot be empty")
		return nil
	}

	_, ok := c.App.Metadata[username]
	if ok {
		fmt.Fprintf(os.Stdin, "Account with username %s already exists.\n", username)
		return nil
	}

	c.App.Metadata[username] = &account{
		transactionPin: txPin,
	}

	fmt.Fprintf(os.Stdin, "Account with username %s has been created successfully.\n", username)
	return nil
}

// depositToAccount credits an account with the amount specified via the
// command-line.
func depositToAccount(c *cli.Context) error {
	username := c.String("username")
	txPin := c.String("pin")
	amount := c.Float64("amount")

	if amount <= 0 {
		fmt.Fprintln(os.Stderr, "deposit amount must be greater than zero")
		return nil
	}

	account, err := retrieveAndVerifyAccount(c, username, txPin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	account.balance += amount

	fmt.Fprintf(os.Stdin, "Deposit to %s's account was successful. Your new balance is %s\n", username, stringAmount(account.balance, naira))
	return nil
}

// withdrawFromAccount debits the amount specified on the command-line from an
// account.
func withdrawFromAccount(c *cli.Context) error {
	username := c.String("username")
	txPin := c.String("pin")
	amount := c.Float64("amount")

	if amount <= 0 {
		return errors.New("withdrawal amount must be greater than zero")
	}

	account, err := retrieveAndVerifyAccount(c, username, txPin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	if account.balance < amount {
		fmt.Fprintf(os.Stderr, "Sorry, your account balance is insufficient for this withdrawal. You have %s and want to withdraw %s. \n", stringAmount(account.balance, naira), stringAmount(amount, naira))
		return nil
	}
	account.balance -= amount

	fmt.Fprintf(os.Stdin, "Your withdrawal of %s has been concluded successfully. Your remaining balance is: %s\n", stringAmount(amount, naira), stringAmount(account.balance, naira))
	return nil
}

// checkAccountBalance checks and reports an account's balance.
func checkAccountBalance(c *cli.Context) error {
	username := c.String("username")
	txPin := c.String("pin")

	account, err := retrieveAndVerifyAccount(c, username, txPin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	fmt.Fprintf(os.Stdin, "Your account balance is %s.\n", stringAmount(account.balance, naira))
	return nil
}

// changeAccountPin changes an account's transaction pin.
func changeAccountPin(c *cli.Context) error {
	username := c.String("username")
	txPin := c.String("pin")
	newTxPin := c.String("newpin")

	if newTxPin == "" {
		fmt.Fprintln(os.Stderr, "new transaction pin cannot be empty")
		return nil
	}

	userAccount, err := retrieveAndVerifyAccount(c, username, txPin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	userAccount.transactionPin = newTxPin

	fmt.Fprintln(os.Stdin, "Your transaction pin has been changed successfully.")
	return nil
}

// logOut exits the CLI ATM program.
func logOut(c *cli.Context) error {
	fmt.Fprintln(os.Stdin, "Exiting CLI ATM Machine ....")
	os.Exit(0)
	return nil
}

// retrieveAndVerifyAccount retrieves and verify a user account.
func retrieveAndVerifyAccount(c *cli.Context, username, txPin string) (*account, error) {
	if username == "" || txPin == "" {
		return nil, errors.New("transaction pin or username cannot be empty")
	}

	user, ok := c.App.Metadata[username]
	if !ok {
		return nil, fmt.Errorf("user with username %s does not exist. Please proceed to create an account", username)
	}

	userAccount := user.(*account)
	if userAccount.transactionPin != txPin {
		return nil, errors.New("unauthorized: provided transaction pin does not match account's transaction pin")
	}
	return userAccount, nil
}

// performOtherTransaction prompts the user to perform other transactions.
func performOtherTransaction(c *cli.Context) error {
	fmt.Println("Do you wish to perform other transactions? \nIf yes, provide the required arguments in the format below:\ndeposit ---amount=2000 --username=lemon --pin=1234")
	fmt.Println("If you don't wish to continue, exit the app with: logout")

	args := ReadArgs()
	if len(args) < 1 {
		fmt.Fprintln(c.App.ErrWriter, "unexpected arguments provided expected more one or more argument")
		return nil
	}

	return c.App.RunContext(c.Context, args)
}

// stringAmount returns an amount with it's currency formatted as a string.
func stringAmount(amount float64, currency string) string {
	return fmt.Sprintf("%.2f %s", amount, currency)
}

// ReadArgs retrieves user input from the command-line.
func ReadArgs() []string {
	args := []string{"cli-atm"}
	scanner := bufio.NewScanner(os.Stdin)
	var textInput string
	for scanner.Scan() {
		if textInput = scanner.Text(); textInput != "" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error reading standard input: %v", err))
		return nil
	}

	newArgs := make([]string, 0)
	for _, arg := range strings.SplitAfter(textInput, " ") {
		args = append(args, strings.Trim(arg, " "))
	}

	args = append(args, newArgs...)
	return args
}
