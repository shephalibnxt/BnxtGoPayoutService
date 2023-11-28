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

type FundAccountService struct {
	Database database.FundAccountDatabase
}

func (fas FundAccountService) RazorpayApiCallForFundAccount(c *gin.Context, bankAccount, ifsc, name string, requestFundAccountPayload entity.RequestToFundAccountApi) (map[string]interface{}, error) {
	configInstance := config.InitConfig()
	rp_username := configInstance.String("rp_username")
	rp_password := configInstance.String("rp_password")
	razorpay_fundaccount_url := configInstance.String("razorpay_fundaccount_url")

	
	finalRequest := entity.RequestFuncAccount{
		ContactId:   requestFundAccountPayload.ContactId,
		AccountType: "bank_account",
		BankAccount: struct {
			Name          string `json:"name"`
			Ifsc          string `json:"ifsc"`
			AccountNumber string `json:"account_number"`
		}{
			Name:          name,
			Ifsc:          ifsc,
			AccountNumber: bankAccount,
		},
	}

	// Serialize the 'requestFundAccountPayload' struct into a JSON byte slice
	fmt.Println("payload before:", finalRequest)
	requestPayloadBytes, err := json.Marshal(finalRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}
	fmt.Println("payload after:", finalRequest)

	// Create an HTTP request with basic authentication in the header
	request, err := http.NewRequest("POST", razorpay_fundaccount_url, bytes.NewBuffer(requestPayloadBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(rp_username, rp_password)

	//send the http request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	//close response body
	defer response.Body.Close()

	//read all response data
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	//Deserialize the JSON response from the Razorpay API, contained in 'responseBody', into a Go map 'RazorpayResponse'.
	var RazorpayResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &RazorpayResponse)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	return RazorpayResponse, nil
}
