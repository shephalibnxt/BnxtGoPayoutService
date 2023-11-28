package database

import (
	"database/sql"
	"event-service/entity"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PayoutBankAccountDatabase struct {
	Db *sql.DB
}

// retrieve fundAccountId
func (pbad PayoutBankAccountDatabase) RetrieveFundAccountId(c *gin.Context, createPayoutRequest entity.RequestToPayoutAccount) (string, error) {

	var fundAccountId string
	fundIDSelectStatement := `select fund_account_id from bharat_nxt_payment.razorpay_contact_details
	where payee_id = $1`

	err := pbad.Db.QueryRow(fundIDSelectStatement, createPayoutRequest.PayeeId).Scan(&fundAccountId)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}

	fmt.Println("fund_account_ids: ", fundAccountId)
	return fundAccountId, nil
}

// retrieve reference id
func (pbad PayoutBankAccountDatabase) RetrieveReferenceId(c *gin.Context, createPayoutRequest entity.RequestToPayoutAccount) (string, error) {

	var referenceId string
	referenceIdSelectStatement := `select order_id from bharat_nxt_payment.payment_details 
	where id = $1`

	err := pbad.Db.QueryRow(referenceIdSelectStatement, createPayoutRequest.PaymentId).Scan(&referenceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}

	fmt.Println("order_id: ", referenceId)
	return referenceId, nil
}

// retrieve razorpay contact details id
func (pbad PayoutBankAccountDatabase) RetrieveRazorpayContactDetailsId(c *gin.Context, FundAccountId string) (razorpayContactDetailsId int64) {
	//var razorpayContactDetailsId int64
	razorpayContactDetailsIdSelectStatement := `select id  from bharat_nxt_payment.razorpay_contact_details 
	where fund_account_id = $1`

	err := pbad.Db.QueryRow(razorpayContactDetailsIdSelectStatement, FundAccountId).Scan(&razorpayContactDetailsId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Razorpay contact details id:", razorpayContactDetailsId)

	return razorpayContactDetailsId
}

// retrieve amount
func (pbad PayoutBankAccountDatabase) RetrieveAmount(c *gin.Context, createPayoutRequest entity.RequestToPayoutAccount) (int64, error) {
	var amount int64
	amountJoinStatement := `SELECT td.payee_amount FROM bharat_nxt_payment.transaction_details td
	LEFT JOIN bharat_nxt_payment.payment_details pd ON pd.bnxt_txn_id = pd.id
	LEFT JOIN bharat_nxt_payment.razorpay_contact_details rcd ON pd.payee_id::varchar = rcd.contact_id
	LEFT JOIN bharat_nxt_payment.payee_details pyd ON rcd.payee_id::varchar = pyd.id::varchar
	WHERE pyd.id = $1;`
	err := pbad.Db.QueryRow(amountJoinStatement, createPayoutRequest.PayeeId).Scan(&amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 0.0, err

	}

	fmt.Println("Amount: ", amount)
	return amount, nil
}
func (pbad PayoutBankAccountDatabase) SavePayoutDetails(c *gin.Context, razorpayResponseData map[string]interface{}, RazorpayContactDetailsId int64) {
	insertStatement := `INSERT INTO bharat_nxt_payment.payout_details (status,
		batch_id, failure_reason, amount, utr, created_at,updated_at,payout_id, razorpay_contact_details_id) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	currentTimeStamp := time.Now()

	formattedTimeStamp := currentTimeStamp.Format("2006-01-02 15:04:05")

	_, err := pbad.Db.Exec(insertStatement, razorpayResponseData["status"], razorpayResponseData["batch_id"],
		razorpayResponseData["failure_reason"], razorpayResponseData["amount"], razorpayResponseData["utr"], formattedTimeStamp,
		formattedTimeStamp, razorpayResponseData["id"], RazorpayContactDetailsId)
	if err != nil {
		errorMessage := fmt.Sprintf("Database insertion error: %v", err)
		fmt.Println(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
