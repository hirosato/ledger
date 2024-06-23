package ledger

type Account struct {
	Name     string
	Balance  map[string]float64
	Children map[string]*Account
}

func NewAccount(name string) *Account {
	return &Account{
		Name:     name,
		Balance:  make(map[string]float64),
		Children: make(map[string]*Account),
	}
}

func (a *Account) AddAmount(amount float64, currency string) {
	a.Balance[currency] += amount
}

func (a *Account) AddChild(name string) *Account {
	if child, exists := a.Children[name]; exists {
		return child
	}
	child := NewAccount(name)
	a.Children[name] = child
	return child
}
