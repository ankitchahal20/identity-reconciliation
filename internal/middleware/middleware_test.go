package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/identity-reconciliation/internal/constants"
	"github.com/identity-reconciliation/internal/models"
	"github.com/identity-reconciliation/internal/utils"
	"gotest.tools/assert"
)


func TestValidateRequestInput(t *testing.T) {
	// init logging client
	utils.InitLogClient()

	// Case 1 : email missing
	requestFields := models.ContactRequest{
		PhoneNumber: "1234567890",
	}

	jsonValue, _ := json.Marshal(requestFields)

	w := httptest.NewRecorder()
	_, e := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/v1/identify", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Case 2 : invalid email
	requestFields = models.ContactRequest{
		PhoneNumber: "1234567890",
		Email: "ankitchahalgmail.com",
		
	}

	jsonValue, _ = json.Marshal(requestFields)

	w = httptest.NewRecorder()
	_, e = gin.CreateTestContext(w)
	req, _ = http.NewRequest(http.MethodPost, "/v1/identify", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Case 3 : phone number missing 
	requestFields = models.ContactRequest{
		Email: "ankitchahal@gmail.com",
	}

	jsonValue, _ = json.Marshal(requestFields)

	w = httptest.NewRecorder()
	_, e = gin.CreateTestContext(w)
	req, _ = http.NewRequest(http.MethodPost, "/v1/identify", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}