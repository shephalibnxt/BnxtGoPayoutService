package entity

type PayoutEvent struct {
	Entity    string   `json:"entity"`
	AccountId string   `json:"account_id"`
	Event     string   `json:"event"`
	Contains  []string `json:"contains"`
	Payload   struct {
		Payout struct {
			Entity struct {
				Id            string            `json:"id"`
				Entity        string            `json:"entity"`
				FundAccountId string            `json:"fund_account_id"`
				Amount        int64             `json:"amount"`
				Currency      string            `json:"currency"`
				Notes         map[string]string `json:"notes"`
				Fees          int64             `json:"fees"`
				Tax           int64             `json:"tax"`
				Status        string            `json:"status"`
				Purpose       string            `json:"purpose"`
				UTR           string            `json:"utr"`
				Mode          string            `json:"mode"`
				ReferenceId   string            `json:"reference_id"`
				Narration     string            `json:"narration"`
				BatchId       string            `json:"batch_id"`
				StatusDetails struct {
					Description string `json:"description"`
					Source      string `json:"source"`
					Reason      string `json:"reason"`
				} `json:"status_details"`
				CreatedAt int64  `json:"created_at"`
				FeeType   string `json:"fee_type"`
			} `json:"entity"`
		} `json:"payout"`
	} `json:"payload"`
	CreatedAt int64 `json:"created_at"`
}
