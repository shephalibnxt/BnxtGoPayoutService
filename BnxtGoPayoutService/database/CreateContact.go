package database

import (
	"database/sql"
	"event-service/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactDatabase struct {
	Db *sql.DB
}

// func (cd ContactDatabase) GetContactDetails(c *gin.Context, requestPayload entity.RequestToContactApi) (Name, Email, Phone string) {

// 	//check if user already exists or not
// 	var contactId string

// 	existingContactQuery := `SELECT contact_id FROM bharat_nxt_payment.razorpay_contact_details
// 	WHERE payee_id = $1;`

// 	err := cd.Db.QueryRow(existingContactQuery, requestPayload.PayeeId).Scan(&contactId)
// 	if err == nil {
// 		c.JSON(http.StatusOK, gin.H{"contact_id": contactId})
// 		return
// 	} else {
// 		//retrieve name, email, phone from payee_id and user id

// 		nameSelectStatement := `select pd.name , upm.payee_email , ROUND(upm.payee_phone)::VARCHAR
// 			from bharat_nxt_payment.user_payee_mapping upm
// 			join bharat_nxt_payment.payee_details pd on upm.payee_id = pd.id
// 			where upm.user_id = $1 and upm.payee_id = $2;`
// 		err := cd.Db.QueryRow(nameSelectStatement, requestPayload.UserId, requestPayload.PayeeId).Scan(&Name, &Email, &Phone)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Values does not exist in database."})
// 			return
// 		}
// 		fmt.Println("Name: ", Name)
// 		fmt.Println("Email: ", Email)
// 		fmt.Println("Phone: ", Phone)

// 		// Return all the fields in the response
// 		// c.JSON(http.StatusOK, gin.H{
// 		// 	"contact_id": contactId,
// 		// 	"name":       Name,
// 		// 	"email":      Email,
// 		// 	"phone":      Phone,
// 		// })

// 		return Name, Email, Phone
// 	}

// }

func (cd ContactDatabase) GetContactID(c *gin.Context, requestPayload entity.RequestToContactApi) (contactId string) {
	//check if user already exists or not
	//var contactId string

	existingContactQuery := `SELECT contact_id FROM bharat_nxt_payment.razorpay_contact_details
	WHERE payee_id = $1;`

	err := cd.Db.QueryRow(existingContactQuery, requestPayload.PayeeId).Scan(&contactId)
	if err == nil {

		return
	}

	return contactId
}

func (cd ContactDatabase) RetrieveUserDetails(c *gin.Context, requestPayload entity.RequestToContactApi) (Name, Email, Phone string) {
	//retrieve name, email, phone from payee_id and user id
	nameSelectStatement := `select pd.name , upm.payee_email , ROUND(upm.payee_phone)::VARCHAR
		from bharat_nxt_payment.user_payee_mapping upm
		join bharat_nxt_payment.payee_details pd on upm.payee_id = pd.id
		where upm.user_id = $1 and upm.payee_id = $2;`

	err := cd.Db.QueryRow(nameSelectStatement, requestPayload.UserId, requestPayload.PayeeId).Scan(&Name, &Email, &Phone)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	fmt.Printf("error:%v", err)

	fmt.Println("Name: ", Name)
	fmt.Println("Email: ", Email)
	fmt.Println("Phone: ", Phone)

	return Name, Email, Phone
}

func (cd ContactDatabase) SaveContactDetails(c *gin.Context, RazorpayResponse map[string]interface{}, requestPayload entity.RequestToContactApi) {

	//insert fields from response into postgresSQL database
	insertStatement := `INSERT INTO bharat_nxt_payment.razorpay_contact_details (contact_id,is_active,payee_id,fund_account_id)
	VALUES($1,$2,$3::int,$4)`

	_, err := cd.Db.Exec(insertStatement, RazorpayResponse["id"], RazorpayResponse["active"], requestPayload.PayeeId, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
