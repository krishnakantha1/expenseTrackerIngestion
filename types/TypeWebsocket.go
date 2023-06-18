package types

type ExpenseMessage struct {
	ExpsenseEntry ExpsenseEntry `json:"expenseEntry"`
	RawMessage    RawMessage    `json:"rawMessage"`
}

type ExpsenseEntry struct {
	URI             string `json:"uri"`
	Bank            string `json:"bank"`
	EncryptedAmount string `json:"encryptedAmount"`
	ExpenseDate     string `json:"expenseData_long"`
	ExpenseType     string `json:"expenseType"`
	ExpenseTag      string `json:"tag"`
}

type RawMessage struct {
	Raw string `json:"raw"`
}
