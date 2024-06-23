package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hirosato/ledger"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ledger",
	Short: "Ledger is a double-entry accounting system",
	Long:  `A modern implementation of ledger-cli in Go, using YAML for data storage.`,
}

var balanceCmd = &cobra.Command{
	Use:   "balance [account]",
	Short: "Show account balance",
	Long:  `Show the balance of a specific account or all accounts if no account is specified`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		l := getLedger()
		var account *ledger.Account

		if len(args) == 0 {
			account = l.Balance("") // Get all accounts
		} else {
			account = l.Balance(args[0])
		}

		printAccountHierarchy(account, "", make(map[string]float64))
	},
}

var registerCmd = &cobra.Command{
	Use:   "register [account]",
	Short: "Show account transactions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		l := getLedger()
		transactions := l.Register(args[0])
		printRegister(transactions, args[0], l)
	},
}

var filePath string
var defaultCurrency string

func init() {
	RootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "ledger.yaml", "ledger file (default is ./ledger.yaml)")
	RootCmd.PersistentFlags().StringVarP(&defaultCurrency, "currency", "c", "", "default currency")
	RootCmd.AddCommand(balanceCmd)
	RootCmd.AddCommand(registerCmd)
}

func Execute() error {
	return RootCmd.Execute()
}

func getLedger() *ledger.Ledger {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	l, err := ledger.ParseYAML(data)
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	return l
}

func printAccountHierarchy(account *ledger.Account, indent string, totals map[string]float64) {
	if account.Name != "" {
		fmt.Printf("%s%s\n", indent, account.Name)
		for currency, amount := range account.Balance {
			fmt.Printf("%s  %s %.2f\n", indent, currency, amount)
			totals[currency] += amount
		}
	}

	var childNames []string
	for name := range account.Children {
		childNames = append(childNames, name)
	}
	sort.Strings(childNames)
	for _, name := range childNames {
		printAccountHierarchy(account.Children[name], indent+"  ", totals)
	}

	if indent == "" {
		fmt.Println("\nTotal:")
		for currency, amount := range totals {
			fmt.Printf("  %s %.2f\n", currency, amount)
		}
	}
}

func printRegister(transactions []ledger.Transaction, accountName string, l *ledger.Ledger) {
	totals := make(map[string]float64)
	for _, t := range transactions {
		fmt.Printf("%s %s\n", t.Date.Format("2006-01-02"), t.Description)
		for _, p := range t.Postings {
			if strings.HasPrefix(p.Account, accountName) {
				currency := l.GetCurrency(p)
				fmt.Printf("  %s  %.2f %s\n", p.Account, p.Amount, currency)
				totals[currency] += p.Amount
				fmt.Printf("    Total: %.2f %s\n", totals[currency], currency)
			}
		}
	}

	fmt.Println("\nFinal Totals:")
	for currency, amount := range totals {
		fmt.Printf("  %s %.2f\n", currency, amount)
	}
}
