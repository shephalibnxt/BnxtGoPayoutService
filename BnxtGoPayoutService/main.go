package main

import (
	"event-service/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {
	var Data entity.PayoutEvent

	err := c.BindJSON(&Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	fmt.Printf("Received event data:\n")
	fmt.Printf("Entity: %s\n", Data.Entity)
	fmt.Printf("Account ID: %s\n", Data.AccountId)
	fmt.Printf("Event: %s\n", Data.Event)
	fmt.Printf("Contains: %v\n", Data.Contains)
	fmt.Printf("Payout Id: %s\n", Data.Payload.Payout.Entity.Id)
	fmt.Printf("Payout Entity: %s\n", Data.Payload.Payout.Entity.Entity)
	fmt.Printf("Fund Account Id: %s\n", Data.Payload.Payout.Entity.FundAccountId)
	fmt.Printf("Amount: %d\n", Data.Payload.Payout.Entity.Amount)
	fmt.Printf("Currency: %s\n", Data.Payload.Payout.Entity.Currency)
	fmt.Printf("Notes: %+v\n", Data.Payload.Payout.Entity.Notes)
	fmt.Printf("Fees: %d\n", Data.Payload.Payout.Entity.Fees)
	fmt.Printf("Tax: %d\n", Data.Payload.Payout.Entity.Tax)
	fmt.Printf("Status: %s\n", Data.Payload.Payout.Entity.Status)
	fmt.Printf("UTR: %s\n", Data.Payload.Payout.Entity.UTR)
	fmt.Printf("Mode: %s\n", Data.Payload.Payout.Entity.Mode)
	fmt.Printf("Reference Id: %s\n", Data.Payload.Payout.Entity.ReferenceId)
	fmt.Printf("Narration: %s\n", Data.Payload.Payout.Entity.Narration)
	fmt.Printf("Batch Id: %s\n", Data.Payload.Payout.Entity.BatchId)
	fmt.Printf("Status Details:\n")
	fmt.Printf("Description: %s\n", Data.Payload.Payout.Entity.StatusDetails.Description)
	fmt.Printf("Source: %s\n", Data.Payload.Payout.Entity.StatusDetails.Source)
	fmt.Printf("Reason: %s\n", Data.Payload.Payout.Entity.StatusDetails.Reason)
	fmt.Printf("Fee Type: %s\n", Data.Payload.Payout.Entity.FeeType)
	fmt.Printf("Created At: %d\n", Data.CreatedAt)

	c.JSON(http.StatusOK, Data)
}

func main() {
	router := gin.Default()

	router.POST("/pendingEvent", CreateEvent)
	fmt.Println("Server is running...")
	router.Run(":8080")
}
