package middleware

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/identity-reconciliation/internal/constants"
	"github.com/identity-reconciliation/internal/models"
	"github.com/identity-reconciliation/internal/utils"
)

// This function gets the unique transactionID
func getTransactionID(c *gin.Context) string {
	transactionID := c.GetHeader(constants.TransactionID)
	_, err := uuid.Parse(transactionID)
	if err != nil {
		transactionID = uuid.New().String()
		c.Set(constants.TransactionID, transactionID)
	}
	return transactionID
}

func ValidateInputRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get the transactionID from headers if not present create a new.
		transactionID := getTransactionID(ctx)
		validateIdentifyInput(ctx, transactionID)
		ctx.Next()
	}
}

func validateIdentifyInput(ctx *gin.Context, txid string) {
	utils.Logger.Info(fmt.Sprintf("request received for identify endpoint, txid : %v", txid))
	var contact models.ContactRequest
	err := ctx.ShouldBindBodyWith(&contact, binding.JSON)
	if err != nil {
		utils.Logger.Error("error while unmarshaling the request field for identify data validation")
		utils.RespondWithError(ctx, http.StatusInternalServerError, constants.InvalidIdentifyBody)
		return
	}

	if contact.Email == "" {
		utils.Logger.Error(fmt.Sprintf("email field is missing, txid : %v", txid))
		errMessage := "email field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}

	_, parseErr := mail.ParseAddress(contact.Email)
		if parseErr != nil {
			utils.Logger.Error(fmt.Sprintf("email received is incorrect, txid : %v", txid))
			err := fmt.Errorf("invalid email found, err : %v", parseErr)
			utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
			return
		}

	if contact.PhoneNumber == "" {
		utils.Logger.Error(fmt.Sprintf("phone_number field is missing, txid : %v", txid))
		errMessage := "phone_number field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}

	// add check for phone number validation
	// pattern := `^(0|91|\+91)?[6-9]\d{9}$`
	// re := regexp.MustCompile(pattern)
	// phoneNumber := strings.ReplaceAll(contact.PhoneNumber, " ", "")
    // if re.Find([]byte(phoneNumber)) == nil {
	// 	utils.Logger.Error(fmt.Sprintf("phone_number field is invalid, txid : %v", txid))
	// 	errMessage := "phone_number field is invalid"
	// 	utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
	// 	return
	// }
}
