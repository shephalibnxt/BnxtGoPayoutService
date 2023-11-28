package service

import (
	"bytes"
	"encoding/json"
	"event-service/config"
	"event-service/database"
	"event-service/entity"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PayoutBankAccountService struct {
	Database database.PayoutBankAccountDatabase
}

func (pbas PayoutBankAccountService) RazorpayApiCallForPayout(c *gin.Context, FundAccountId, ReferenceId string, RazorpayContactDetailsId, Amount int64) (map[string]interface{}, error) {
	configInstance := config.InitConfig()
	rp_username := configInstance.String("rp_username")
	rp_password := configInstance.String("rp_password")
	razorpay_payout_url := configInstance.String("razorpay_payout_url")

	razorpayRequest := entity.RequestPayoutBankAccount{
		AccountNumber:     "2323230073798533",
		FundAccountId:     FundAccountId,
		Amount:            Amount,
		Currency:          "INR",
		Mode:              "IMPS",
		Purpose:           "refund",
		QueueIfLowBalance: true,
		ReferenceId:       ReferenceId,
		Narration:         "Acme Corp Fund Transfer",
		Notes: map[string]string{
			"notes_key_1": "Tea, Earl Grey, Hot",
			"notes_key_2": "Tea, Earl Grey… decaf.",
		},
	}

	fmt.Println("Final request.......:", razorpayRequest)

	razorpayRequestBytes, err := json.Marshal(razorpayRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	request, err := http.NewRequest("POST", razorpay_payout_url, bytes.NewBuffer(razorpayRequestBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err

	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(rp_username, rp_password)

	client := &http.Client{}
	razorpayResponse, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err

	}

	defer razorpayResponse.Body.Close()

	responseBody, err := io.ReadAll(razorpayResponse.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return nil, err

	}

	var razorpayResponseData map[string]interface{}
	err = json.Unmarshal(responseBody, &razorpayResponseData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err

	}
	fmt.Println("response data...........:", razorpayResponseData)

	return razorpayResponseData, nil
}
