package main

import (
	"log"

	"github.com/identity-reconciliation/internal/config"
	"github.com/identity-reconciliation/internal/db"
	"github.com/identity-reconciliation/internal/server"
	"github.com/identity-reconciliation/internal/service"
	"github.com/identity-reconciliation/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Initializing the Log client
	utils.InitLogClient()

	// Initializing the GlobalConfig
	err := config.InitGlobalConfig()
	if err != nil {
		log.Fatalf("Unable to initialize global config")
	}

	// Establishing the connection to DB.
	postgres, err := db.New()
	if err != nil {
		log.Fatal("Unable to connect to DB : ", err)
	}

	// Initializing the client for notes service
	_ = service.NewIdentityReconciliationService(postgres)

	// Starting the server
	server.Start()
}
