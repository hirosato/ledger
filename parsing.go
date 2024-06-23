package ledger

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func ParseYAML(data []byte) (*Ledger, error) {
	var rawLedger struct {
		Transactions []struct {
			Date        string `yaml:"date"`
			Description string `yaml:"description"`
			Postings    []struct {
				Account string `yaml:"account"`
				Amount  string `yaml:"amount"`
			} `yaml:"postings"`
		} `yaml:"transactions"`
	}

	err := yaml.Unmarshal(data, &rawLedger)
	if err != nil {
		return nil, err
	}

	ledger := New()
	for _, rawTxn := range rawLedger.Transactions {
		txn := Transaction{
			Date:        parseDate(rawTxn.Date),
			Description: rawTxn.Description,
		}

		for _, rawPosting := range rawTxn.Postings {
			posting, err := ParsePosting(rawPosting.Account, rawPosting.Amount)
			if err != nil {
				return nil, err
			}
			txn.Postings = append(txn.Postings, posting)
		}

		ledger.AddTransaction(txn)
	}

	return ledger, nil
}

func ParsePosting(account, amountStr string) (Posting, error) {
	re := regexp.MustCompile(`^([A-Z]{3})?\s*(-?\d+(\.\d+)?)`)
	matches := re.FindStringSubmatch(strings.TrimSpace(amountStr))

	if matches == nil {
		return Posting{}, errors.New("invalid amount format")
	}

	currency := matches[1] // This will be empty if no currency is specified
	amountVal, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return Posting{}, err
	}

	return Posting{
		Account:  account,
		Amount:   amountVal,
		Currency: currency,
	}, nil
}

func parseDate(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{} // Return zero time if parsing fails
	}
	return t
}
