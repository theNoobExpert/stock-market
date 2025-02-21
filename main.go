package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/theNoobExpert/icicibreeze/connect"
	"github.com/theNoobExpert/icicibreeze/pkg/utils"
)

var logger = utils.GetLogger()

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Could not load env file. Check if .env file exists.")
		panic("unable to read .env file.")
	}
	logger.Debug("Env loaded successfully")
	timeoutSecs := 300
	timeoutStr, exists := os.LookupEnv("BREEZE_HTTP_TIMEOUT_SECS")
	if exists {
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			logger.Warnf("Error while setting timeout from env. Using default timeout: %d. Error: %w", timeoutSecs, err)
		} else {
			timeoutSecs = timeout
		}
	} else {
		logger.Warn(fmt.Sprintf("BREEZE_HTTP_TIMEOUT_SECS env not set. Setting default http timeout of %d.", timeoutSecs))
	}
	breezeClient, err := connect.NewBreezeConnectClient(os.Getenv("BREEZE_APP_KEY"), os.Getenv("BREEZE_APP_SECRET"), "", timeoutSecs)
	if err != nil {
		logger.Fatalf("Error while creating breeze client: %w", err)
		return
	}
	loginUrl, _ := breezeClient.GetLoginURL()
	logger.Info(loginUrl)

	customerDetails, err := breezeClient.InitSessionToken(os.Getenv("BREEZE_SESSION_TOKEN"))
	if err != nil {
		logger.Fatalf("error while getting session token: %w", err)
	}

	jsonData, err := json.MarshalIndent(customerDetails, "", "  ")
	if err != nil {
		logger.Errorf("Error while marshalling customer funds: %v", err)
		return
	}
	logger.Infof("Customer details : \n%s", jsonData)
}
