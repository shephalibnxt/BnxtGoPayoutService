package database

import (
	"database/sql"
	"event-service/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FundAccountDatabase struct {
	Db *sql.DB
}

func (fad FundAccountDatabase) GetAccountDetails(c *gin.Context, requestFundAccountPayload entity.RequestToFundAccountApi) (BankAccount, Ifsc, Name string) {
	//retrieve bank account number, ifsc, name from contact_id
	//var BankAccount, Ifsc, Name string
	joinStatement := `SELECT  pd.bank_ac, pd.ifsc, pd.name
	FROM bharat_nxt_payment.payee_details pd
	JOIN bharat_nxt_payment.razorpay_contact_details rcd ON pd.id::character varying = rcd.payee_id
	WHERE rcd.contact_id = $1;`

	err := fad.Db.QueryRow(joinStatement, requestFundAccountPayload.ContactId).Scan(&BankAccount, &Ifsc, &Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Values does not exist in database."})
		return
	}

	fmt.Println(" bankAccount: ", BankAccount)
	fmt.Println(" ifsc: ", Ifsc)
	fmt.Println(" name: ", Name)
	return BankAccount, Ifsc, Name

}

func (fad FundAccountDatabase) SaveFundAccountDetails(c *gin.Context, RazorpayResponse map[string]interface{}) {
	selectQuery := "SELECT * from bharat_nxt_payment.razorpay_contact_details WHERE contact_id = $1"
	result, err := fad.Db.Query(selectQuery, RazorpayResponse["contact_id"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !result.Next() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact Id does not exist in the table"})
	} else {
		fmt.Println("Contact id found ", result)
		updateQuery := "UPDATE bharat_nxt_payment.razorpay_contact_details SET fund_account_id = $1 where contact_id = $2"
		_, err := fad.Db.Exec(updateQuery, RazorpayResponse["id"], RazorpayResponse["contact_id"])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

}
