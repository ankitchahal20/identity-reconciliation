package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/identity-reconciliation/internal/constants"
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
		fmt.Printf("TimeStamp : %v", time.Now().UTC, transactionID)

		ctx.Next()
	}
}
