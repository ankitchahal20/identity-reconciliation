package service

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/identity-reconciliation/internal/db"
	"github.com/identity-reconciliation/internal/utils"
)

var (
	identityReconciliationClient *IdentityReconciliationService
	once                         sync.Once
)

type IdentityReconciliationService struct {
	repo db.IdentityReconciliationService
}

// creditCardLimitOfferClient should only be created once throughtout the application lifetime
func NewIdentityReconciliationService(conn db.IdentityReconciliationService) *IdentityReconciliationService {
	if identityReconciliationClient == nil {
		once.Do(
			func() {
				identityReconciliationClient = &IdentityReconciliationService{
					repo: conn,
				}
			})
	} else {
		utils.Logger.Info("identityReconciliationClient is alredy created")
	}
	return identityReconciliationClient
}

// This function is responsible for account creation
func Identify() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

	}
}
