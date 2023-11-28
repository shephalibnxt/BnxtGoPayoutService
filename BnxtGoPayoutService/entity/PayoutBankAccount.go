package entity

type RequestToPayoutAccount struct {
	PaymentId string `json:"payment_id"`
	PayeeId   string `json:"payee_id"`
}

// request struct
type RequestPayoutBankAccount struct {
	AccountNumber     string            `json:"account_number"`
	FundAccountId     string            `json:"fund_account_id"`
	Amount            int64             `json:"amount"`
	Currency          string            `json:"currency"`
	Mode              string            `json:"mode"`
	Purpose           string            `json:"purpose"`
	QueueIfLowBalance bool              `json:"queue_if_low_balance"`
	ReferenceId       string            `json:"reference_id"`
	Narration         string            `json:"narration"`
	Notes             map[string]string `json:"notes"`
}

// Response struct
type ResponsePayoutBankAccount struct {
	ID            string            `json:"id"`
	Entity        string            `json:"entity"`
	FundAccountID string            `json:"fund_account_id"`
	Amount        int64             `json:"amount"`
	Currency      string            `json:"currency"`
	Notes         map[string]string `json:"notes"`
	Fees          int               `json:"fees"`
	Tax           int               `json:"tax"`
	Status        string            `json:"status"`
	Purpose       string            `json:"purpose"`
	UTR           string            `json:"utr"`
	Mode          string            `json:"mode"`
	ReferenceID   string            `json:"reference_id"`
	Narration     string            `json:"narration"`
	BatchID       string            `json:"batch_id"`
	FailureReason string            `json:"failure_reason"`
	CreatedAt     int64             `json:"created_at"`
	FeeType       string            `json:"fee_type"`
	StatusDetails struct {
		Reason      string `json:"reason"`
		Description string `json:"description"`
		Source      string `json:"source"`
	} `json:"status_details"`
	MerchantID      string `json:"merchant_id"`
	StatusDetailsID string `json:"status_details_id"`
	Error           struct {
		Source      string                 `json:"source"`
		Reason      string                 `json:"reason"`
		Description string                 `json:"description"`
		Code        string                 `json:"code"`
		Step        string                 `json:"step"`
		Metadata    map[string]interface{} `json:"metadata"`
	} `json:"error"`
}
