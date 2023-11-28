package main

import (
	"event-service/config"
	"event-service/controller"
	"event-service/database"
	_ "event-service/docs"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Payout Service
// @version 1.0
// @description Payout Service using gin framework
// @host localhost:8080
func main() {
	config := config.InitConfig()
	if config == nil {
		// if any error occurs program will not run
		return
	}

	portnumber := config.String("portNumber")
	router := gin.Default()

	//Register the swagger handler
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	db := database.PostgresConnection()

	ContactControllerObject := controller.ContactController{
		Database: database.ContactDatabase{Db: db},
	}

	FundAccountObject := controller.FundAccountController{
		Database: database.FundAccountDatabase{Db: db},
	}

	PayoutBankAccountObject := controller.PayoutBankAccountController{

		Database: database.PayoutBankAccountDatabase{Db: db},
	}

	router.POST("/pendingEvent", controller.CreateEvent)
	router.POST("bnxt/createContact", ContactControllerObject.CreateContact)
	router.POST("bnxt/createFundAccount", FundAccountObject.CreateFundAccount)
	router.POST("bnxt/createPayoutAccount", PayoutBankAccountObject.CreatePayout)
	fmt.Println("Server is running...")
	router.Run(":" + portnumber)
}
