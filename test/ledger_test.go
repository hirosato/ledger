package test

import (
	"os"
	"testing"
	"time"

	"github.com/hirosato/ledger"
)

func TestLedger(t *testing.T) {
	// Read the test1 YAML file
	data, err := os.ReadFile("test1.yaml")
	if err != nil {
		t.Fatalf("Error reading test1.yaml: %v", err)
	}

	// Parse the YAML data
	l, err := ledger.ParseYAML(data)
	if err != nil {
		t.Fatalf("Error parsing YAML: %v", err)
	}

	// Set a default currency for testing
	l.SetDefaultCurrency("USD")

	// Test the number of transactions
	if len(l.Transactions) != 5 {
		t.Errorf("Expected 5 transactions, got %d", len(l.Transactions))
	}

	// Test the balance of Assets:Checking
	checkingBalance := l.Balance("Assets:Checking")
	expectedBalance := 2424.50
	if checkingBalance.Balance["USD"] != expectedBalance {
		t.Errorf("Expected Assets:Checking balance to be %.2f, got %.2f", expectedBalance, checkingBalance.Balance["USD"])
	}

	// Test the balance of Expenses
	expensesBalance := l.Balance("Expenses")
	expectedExpenses := 1575.50
	if expensesBalance.Balance["USD"] != expectedExpenses {
		t.Errorf("Expected Expenses balance to be %.2f, got %.2f", expectedExpenses, expensesBalance.Balance["USD"])
	}

	// Test the register function for Expenses:Food
	foodTransactions := l.Register("Expenses:Food")
	if len(foodTransactions) != 1 {
		t.Errorf("Expected 1 transaction for Expenses:Food, got %d", len(foodTransactions))
	}
	if foodTransactions[0].Date != time.Date(2023, 6, 3, 0, 0, 0, 0, time.UTC) {
		t.Errorf("Expected transaction date to be 2023-06-03, got %v", foodTransactions[0].Date)
	}
	if foodTransactions[0].Description != "Grocery shopping" {
		t.Errorf("Expected transaction description to be 'Grocery shopping', got '%s'", foodTransactions[0].Description)
	}

	// Test multi-currency support
	jpyBalance := l.Balance("Expenses:Shopping")
	if jpyBalance.Balance["JPY"] != 5000 {
		t.Errorf("Expected JPY balance for Expenses:Shopping to be 5000, got %.2f", jpyBalance.Balance["JPY"])
	}
}
