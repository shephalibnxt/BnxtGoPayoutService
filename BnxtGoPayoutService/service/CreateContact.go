package service

import (
	"bytes"
	"encoding/json"
	"event-service/config"
	"event-service/database"
	"event-service/entity"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactService struct {
	Database database.ContactDatabase
}

func (cc ContactService) RazorpayApiCallForContact(c *gin.Context, name, email, phone string, requestPayload entity.RequestToContactApi) (map[string]interface{}, error) {
	configInstance := config.InitConfig()
	rp_username := configInstance.String("rp_username")
	rp_password := configInstance.String("rp_password")
	razorpay_contact_url := configInstance.String("razorpay_contact_url")

	finalRequest := entity.RequestContact{
		Name:    name,
		Contact: phone,
		Email:   email,
		Type:    "vendor",
		Notes: map[string]string{
			"note_key": "Beam me up Scotty Updated",
		},
	}

	// Serialize the 'finalRequest' struct into a JSON byte slice
	payloadBytes, err := json.Marshal(finalRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}

	// Create an HTTP request with basic authentication in the header
	request, err := http.NewRequest("POST", razorpay_contact_url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(rp_username, rp_password)

	// Send the HTTP request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}
	defer response.Body.Close()

	// Read the entire response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	var razorpayResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &razorpayResponse)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	return razorpayResponse, nil
}
