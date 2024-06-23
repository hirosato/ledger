package test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/hirosato/ledger/cli"
)

func TestCLICommands(t *testing.T) {
	// Save the original args and restore them after the test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name: "Balance Command (All Accounts)",
			args: []string{"ledger", "balance", "-f", "test/test1.yaml", "-c", "USD"},
			expected: []string{
				"Assets",
				"Checking",
				"Savings",
				"CreditCard",
				"Income",
				"Salary",
				"Expenses",
				"Rent",
				"Food",
				"Shopping",
			},
		},
		{
			name: "Balance Command",
			args: []string{"ledger", "balance", "Assets", "-f", "test/test1.yaml", "-c", "USD"},
			expected: []string{
				"Assets",
				"Checking",
				"2424.50",
				"Savings",
				"1000.00",
				"CreditCard",
				"JPY -5000.00",
			},
		},
		{
			name: "Register Command",
			args: []string{"ledger", "register", "Expenses:Food", "-f", "test/test1.yaml", "-c", "USD"},
			expected: []string{
				"2023-06-03 Grocery shopping",
				"Expenses:Food:Groceries  75.50 USD",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdout to capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Set up args for this test
			os.Args = tt.args

			// Run the command
			cli.Execute()

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read the output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Check the output
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', but it didn't.\nGot: %s", expected, output)
				}
			}
		})
	}
}
