package controller

import (
	"bytes"
	"encoding/json"
	"event-service/config"
	"event-service/database"
	"event-service/entity"
	"event-service/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactController struct {
	Database database.ContactDatabase
	Service  service.ContactService
}

// CreateContact godoc
// @Summary Create a new Contact
// @Tags Create Contact API
// @Description Create a new contact based on provided payee id and user id.
// @ID create-contact
// @Accept json
// @Produce json
// @Param requestPayload body entity.RequestToContactApi true "Request payload to create a contact"
// @Success 200 {object} map[string]interface{} "successful response with contact details"
// @Failure 400 {object} map[string]interface{} "Error response when there is a bad request or an api error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /bnxt/createContact [post]
func (cc ContactController) CreateContact(c *gin.Context) {

	configInstance := config.InitConfig()
	rp_username := configInstance.String("rp_username")
	rp_password := configInstance.String("rp_password")

	var requestPayload entity.RequestToContactApi
	err := c.BindJSON(&requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//name, email, phone := cc.Database.GetContactDetails(c, requestPayload)
	contactId := cc.Database.GetContactID(c, requestPayload)
	fmt.Println("contact id", contactId)
	if contactId != "" {
		c.JSON(http.StatusOK, gin.H{"contact_id": contactId})

	} else {

		name, email, phone := cc.Database.RetrieveUserDetails(c, requestPayload)
		// Call the service function
		RazorpayResponse, err := cc.Service.RazorpayApiCallForContact(c, name, email, phone, requestPayload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cc.Database.SaveContactDetails(c, RazorpayResponse, requestPayload)

		sendData := map[string]interface{}{
			"id":      RazorpayResponse["id"],
			"entity":  RazorpayResponse["entity"],
			"name":    RazorpayResponse["name"],
			"contact": RazorpayResponse["contact"],
			"email":   RazorpayResponse["email"],
			"type":    RazorpayResponse["type"],
			//"reference_id": razorpayResponse["reference_id"],
			"batch_id":   RazorpayResponse["batch_id"],
			"active":     RazorpayResponse["active"],
			"notes":      RazorpayResponse["notes"],
			"created_at": RazorpayResponse["created_at"],
		}

		c.JSON(http.StatusOK, gin.H{"message": sendData})

		//var fetchedContactId string
		fetchedContactId := RazorpayResponse["id"].(string)
		payloadForFundAccount := entity.RequestToFundAccount{
			Contactid: fetchedContactId,
		}

		payloadForFundAccountBytes, err := json.Marshal(payloadForFundAccount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fundAccountRequest, err := http.NewRequest("POST", "http://localhost:8080/bnxt/createFundAccount", bytes.NewBuffer(payloadForFundAccountBytes))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fundAccountRequest.Header.Set("Content-Type", "application/json")
		fundAccountRequest.SetBasicAuth(rp_username, rp_password)
		// Send the HTTP request for creating a fund account
		fundAccountClient := &http.Client{}
		fundAccountResponse, err := fundAccountClient.Do(fundAccountRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Close response body for creating a fund account
		defer fundAccountResponse.Body.Close()

	}
}
