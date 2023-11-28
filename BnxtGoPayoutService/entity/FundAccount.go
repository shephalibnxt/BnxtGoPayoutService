package entity

type RequestToFundAccountApi struct {
	ContactId string `json:"contact_id"`
}

// request struct
type RequestFuncAccount struct {
	ContactId   string `json:"contact_id"`
	AccountType string `json:"account_type"`
	BankAccount struct {
		Name          string `json:"name"`
		Ifsc          string `json:"ifsc"`
		AccountNumber string `json:"account_number"`
	} `json:"bank_account"`
}

// response struct
type ResponseFundAccount struct {
	Id          string `json:"id"`
	Entity      string `json:"entity"`
	ContactId   string `json:"contact_id"`
	AccountType string `json:"account_type"`
	BankAccount struct {
		Ifsc          string   `json:"ifsc"`
		BankName      string   `json:"bank_name"`
		Name          string   `json:"name"`
		Notes         []string `json:"notes"`
		AccountNumber string   `json:"account_number"`
	} `json:"bank_account"`
	BatchId   string `json:"batch_id"`
	Active    bool   `json:"active"`
	CreatedAt int64  `json:"created_at"`
}
