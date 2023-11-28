package controller

import (
	"event-service/database"
	"event-service/entity"
	"event-service/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PayoutBankAccountController struct {
	Service  service.PayoutBankAccountService
	Database database.PayoutBankAccountDatabase
}

// CreatePayout godoc
// @Summary Create a new payout
// @Tags Create Payout API
// @Description Create a new payout based on provided payment id and payee id
// @ID create-payout
// @Accept json
// @Produce json
// @Param createPayoutRequest body entity.RequestToPayoutAccount true "Request payload to create a payout"
// @Success 200 {object} map[string]interface{} "successful response with payout details"
// @Failure 400 {object} map[string]interface{} "Error response when there is a bad request or an api"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /bnxt/createPayoutAccount [post]
func (pbac PayoutBankAccountController) CreatePayout(c *gin.Context) {
	var createPayoutRequest entity.RequestToPayoutAccount
	err := c.BindJSON(&createPayoutRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	fmt.Println("Payload:", createPayoutRequest)

	//retrieve fundAccountId
	FundAccountId, err := pbac.Database.RetrieveFundAccountId(c, createPayoutRequest)
	if err != nil {
		fmt.Printf("error:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// retrieve OrderId
	ReferenceId, err := pbac.Database.RetrieveReferenceId(c, createPayoutRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	//retrieve razorpay contact_details_id from fund_account_id
	RazorpayContactDetailsId := pbac.Database.RetrieveRazorpayContactDetailsId(c, FundAccountId)

	//retrieve amount

	Amount, err := pbac.Database.RetrieveAmount(c, createPayoutRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	// fmt.Println("Amount: ", amount)

	razorpayResponseData, err := pbac.Service.RazorpayApiCallForPayout(c, FundAccountId, ReferenceId, RazorpayContactDetailsId, Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	razorpayResponseData = map[string]interface{}{
		"id":              razorpayResponseData["id"],
		"entity":          razorpayResponseData["entity"],
		"fund_account_id": razorpayResponseData["fund_account_id"],
		"amount":          razorpayResponseData["amount"],
		"currency":        razorpayResponseData["currency"],
		"notes":           razorpayResponseData["notes"],
		"fees":            razorpayResponseData["fees"],
		"tax":             razorpayResponseData["tax"],
		"status":          razorpayResponseData["status"],
		"purpose":         razorpayResponseData["purpose"],
		"utr":             razorpayResponseData["utr"],
		"mode":            razorpayResponseData["mode"],
		"reference_id":    razorpayResponseData["reference_id"],
		"narration":       razorpayResponseData["narration"],
		"batch_id":        razorpayResponseData["batch_id"],
		"failure_reason":  razorpayResponseData["failure_reason"],
		"created_at":      razorpayResponseData["created_at"],
		"fee_type":        razorpayResponseData["fee_type"],
		"status_details": map[string]interface{}{
			"reason":      razorpayResponseData["status_details"].(map[string]interface{})["reason"],
			"description": razorpayResponseData["status_details"].(map[string]interface{})["description"],
			"source":      razorpayResponseData["status_details"].(map[string]interface{})["source"],
		},
		"merchant_id":       razorpayResponseData["merchant_id"],
		"status_details_id": razorpayResponseData["status_details_id"],
		"error": map[string]interface{}{
			"source":      razorpayResponseData["error"].(map[string]interface{})["source"],
			"reason":      razorpayResponseData["error"].(map[string]interface{})["reason"],
			"description": razorpayResponseData["error"].(map[string]interface{})["description"],
			"code":        razorpayResponseData["error"].(map[string]interface{})["code"],
			"step":        razorpayResponseData["error"].(map[string]interface{})["step"],
			"metadata":    razorpayResponseData["error"].(map[string]interface{})["metadata"],
		},
	}

	pbac.Database.SavePayoutDetails(c, razorpayResponseData, RazorpayContactDetailsId)

	c.JSON(http.StatusOK, razorpayResponseData)
}
