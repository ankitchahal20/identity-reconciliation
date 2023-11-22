package service

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	identityreconciliationerror "github.com/identity-reconciliation/internal/IdentityReconciliationError"
	"github.com/identity-reconciliation/internal/constants"
	"github.com/identity-reconciliation/internal/db"
	"github.com/identity-reconciliation/internal/models"
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
		txid := ctx.Request.Header.Get(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for identify endpoint, txid : %v", txid))
		var contactInfo models.ContactRequest
		if err := ctx.ShouldBindBodyWith(&contactInfo, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("identify request is unmarshalled successfully, txid : %v", txid))

			reponse, err := identityReconciliationClient.identify(ctx, contactInfo)
			if err != nil {
				utils.RespondWithError(ctx, err.Code, err.Message)
				return
			}

			ctx.JSON(http.StatusOK, reponse)

			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *IdentityReconciliationService) identify(ctx *gin.Context, contactRequest models.ContactRequest) (models.ContactResponse, *identityreconciliationerror.IdentityReconciliationError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	utils.Logger.Info(fmt.Sprintf("calling db layer for identify request, txid : %v", txid))
	createdContact, err := service.repo.FindOrCreateContact(ctx, contactRequest)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer during find/create contact txid : %v", txid))
		return models.ContactResponse{}, err
	}
	fmt.Println("createdContact : ", createdContact)
	return createdContact, nil
}
