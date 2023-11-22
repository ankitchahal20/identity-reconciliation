package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	identityreconciliationerror "github.com/identity-reconciliation/internal/IdentityReconciliationError"
	"github.com/identity-reconciliation/internal/config"
	"github.com/identity-reconciliation/internal/models"
)

var (
	conn *sql.DB
	once sync.Once
)

type postgres struct{ db *sql.DB }

type IdentityReconciliationService interface {
	FindOrCreateContact(*gin.Context, models.ContactRequest) (models.ContactResponse, *identityreconciliationerror.IdentityReconciliationError)
	foundOneRecord(*gin.Context, []models.Contact, string, string) ([]models.Contact, *identityreconciliationerror.IdentityReconciliationError)
	foundMultipleRecord(*gin.Context, []models.Contact, string, string) ([]models.Contact, *identityreconciliationerror.IdentityReconciliationError)
	FindAllContacts(ctx *gin.Context, inputContact models.ContactRequest) ([]models.Contact, *identityreconciliationerror.IdentityReconciliationError)
	handleExistingContact(*gin.Context, []models.Contact, models.ContactRequest) (models.ContactResponse, *identityreconciliationerror.IdentityReconciliationError)
}

func New() (postgres, error) {
	// Initialize the connection only once
	once.Do(func() {
		cfg := config.GetConfig()
		connString := fmt.Sprintf(
			"host=%s dbname=%s password=%s user=%s port=%d",
			cfg.Database.Host, cfg.Database.DBname, cfg.Database.Password,
			cfg.Database.User, cfg.Database.Port,
		)

		var err error
		conn, err = sql.Open("pgx", connString)
		if err != nil {
			log.Fatalf("Unable to connect: %v\n", err)
		}

		log.Println("Connected to database")

		err = conn.Ping()
		if err != nil {
			log.Fatal("Cannot Ping the database")
		}
		log.Println("pinged database")
	})

	return postgres{db: conn}, nil
}
