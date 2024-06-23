package ledger

import (
	"time"
)

type Transaction struct {
	Date        time.Time
	Description string
	Postings    []Posting
}

type Posting struct {
	Account  string
	Amount   float64
	Currency string
}

type TransactionsByDate []Transaction

func (t TransactionsByDate) Len() int           { return len(t) }
func (t TransactionsByDate) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TransactionsByDate) Less(i, j int) bool { return t[i].Date.Before(t[j].Date) }
