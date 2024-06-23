package ledger

import "strings"

type Ledger struct {
	Transactions    []Transaction
	RootAccount     *Account
	DefaultCurrency string
}

func New() *Ledger {
	return &Ledger{
		RootAccount: NewAccount(""),
	}
}

func (l *Ledger) AddTransaction(txn Transaction) {
	l.Transactions = append(l.Transactions, txn)
	for _, posting := range txn.Postings {
		account := l.RootAccount
		parts := strings.Split(posting.Account, ":")
		for _, part := range parts {
			account = account.AddChild(part)
		}
		account.AddAmount(posting.Amount, l.GetCurrency(posting))
	}
}

func (l *Ledger) Balance(accountName string) *Account {
	if accountName == "" {
		return l.RootAccount
	}

	parts := strings.Split(accountName, ":")
	account := l.RootAccount
	for _, part := range parts {
		if part == "" {
			continue
		}
		child, exists := account.Children[part]
		if !exists {
			return NewAccount(part) // Return an empty account if it doesn't exist
		}
		account = child
	}
	return account
}
func (l *Ledger) Register(accountName string) []Transaction {
	var result []Transaction

	for _, txn := range l.Transactions {
		for _, posting := range txn.Postings {
			if strings.HasPrefix(posting.Account, accountName) {
				result = append(result, txn)
				break
			}
		}
	}

	return result
}

func (l *Ledger) SetDefaultCurrency(currency string) {
	l.DefaultCurrency = currency
}

func (l *Ledger) GetCurrency(posting Posting) string {
	if posting.Currency == "" {
		return l.DefaultCurrency
	}
	return posting.Currency
}
