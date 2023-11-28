package entity

type RequestToContactApi struct {
	PayeeId int64 `json:"payee_id"`
	UserId  int64 `json:"user_id"`
}

type RequestToFundAccount struct {
	Contactid string `json:"contact_id"`
}

// request struct
type RequestContact struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Email   string `json:"email"`
	Type    string `json:"type"`
	//ReferenceID string            `json:"reference_id"`
	Notes map[string]string `json:"notes"`
}

// response struct
type ResponseContact struct {
	ID          string            `json:"id"`
	Entity      string            `json:"entity"`
	Name        string            `json:"name"`
	Contact     string            `json:"contact"`
	Email       string            `json:"email"`
	Type        string            `json:"type"`
	ReferenceID string            `json:"reference_id"`
	BatchID     interface{}       `json:"batch_id"`
	Active      bool              `json:"active"`
	Notes       map[string]string `json:"notes"`
	CreatedAt   int64             `json:"created_at"`
}
