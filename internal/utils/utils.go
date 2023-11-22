package utils

import (
	"github.com/gin-gonic/gin"
	identityreconciliationerror "github.com/identity-reconciliation/internal/IdentityReconciliationError"
	"github.com/identity-reconciliation/internal/constants"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogClient() {
	Logger, _ = zap.NewDevelopment()
}

func RespondWithError(c *gin.Context, statusCode int, message string) {

	c.AbortWithStatusJSON(statusCode, identityreconciliationerror.IdentityReconciliationError{
		Trace:   c.Request.Header.Get(constants.TransactionID),
		Code:    statusCode,
		Message: message,
	})
}
