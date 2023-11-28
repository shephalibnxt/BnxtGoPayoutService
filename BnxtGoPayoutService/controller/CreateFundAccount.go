package controller

import (
	"event-service/database"
	"event-service/entity"
	"event-service/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FundAccountController struct {
	Database database.FundAccountDatabase
	Service  service.FundAccountService
}

// CreateFundAccount godoc
// @Summary Create a new Fund Account
// @Tags Create Fund account API
// @Description Create a new fund account based on provided contact id
// @ID create-fund-account
// @Accept json
// @Produce json
// @Param requestFundAccountPayload body entity.RequestToFundAccountApi true "Request payload to create a func account"
// @Success 200 {object} map[string]interface{} "Successful response with fund account details"
// @Failure 400 {object} map[string]interface{} "Error response when there is a bad request or an api error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /bnxt/createFundAccount [post]
func (fac FundAccountController) CreateFundAccount(c *gin.Context) {

	// Parse JSON request body and bind it to the 'requestFundAccountPayload' struct
	var requestFundAccountPayload entity.RequestToFundAccountApi
	err := c.BindJSON(&requestFundAccountPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//retrieve bank account number, ifsc, name from table
	bankAccount, ifsc, name := fac.Database.GetAccountDetails(c, requestFundAccountPayload)

	fmt.Println(" bankAccount: ", bankAccount)
	fmt.Println(" ifsc: ", ifsc)
	fmt.Println(" name: ", name)

	RazorpayResponse, err := fac.Service.RazorpayApiCallForFundAccount(c, bankAccount, ifsc, name, requestFundAccountPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sendData := map[string]interface{}{
		"id":           RazorpayResponse["id"],
		"entity":       RazorpayResponse["entity"],
		"contact_id":   RazorpayResponse["contact_id"],
		"account_type": RazorpayResponse["account_type"],
		"bank_account": map[string]interface{}{
			"ifsc":           RazorpayResponse["bank_account"].(map[string]interface{})["ifsc"],
			"bank_name":      RazorpayResponse["bank_account"].(map[string]interface{})["bank_name"],
			"name":           RazorpayResponse["bank_account"].(map[string]interface{})["name"],
			"notes":          RazorpayResponse["bank_account"].(map[string]interface{})["notes"],
			"account_number": RazorpayResponse["bank_account"].(map[string]interface{})["account_number"],
		},
		"batch_id":   RazorpayResponse["batch_id"],
		"active":     RazorpayResponse["active"],
		"created_at": RazorpayResponse["created_at"],
	}
	fac.Database.SaveFundAccountDetails(c, RazorpayResponse)

	c.JSON(http.StatusOK, gin.H{"message": sendData})
}
