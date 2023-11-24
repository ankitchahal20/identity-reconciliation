package db

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/identity-reconciliation/internal/constants"
	"github.com/identity-reconciliation/internal/models"
	"github.com/identity-reconciliation/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestTransformContact(t *testing.T) {
	// Set up test data
	contactID1 := 1
	contactID2 := 2
	contacts := []models.Contact{
		{
			ID:             &contactID1,
			PhoneNumber:    "123456789",
			Email:          "test@example.com",
			LinkedID:       nil,
			LinkPrecedence: "primary",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             &contactID2,
			PhoneNumber:    "987654321",
			Email:          "another@example.com",
			LinkedID:       new(int),
			LinkPrecedence: "secondary",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	// Call the function to test
	result := transformContact(contacts)

	// Assert that the function behaves as expected
	expected := models.ContactResponse{
		PrimaryContactID:    1,
		Emails:              []string{"test@example.com", "another@example.com"},
		PhoneNumbers:        []string{"123456789", "987654321"},
		SecondaryContactIDs: &[]int{2},
	}
	assert.Equal(t, expected, result)
}

func TestFindOrCreateNewContact(t *testing.T) {
	// mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// instance of the postgres struct with the mock DB
	p := postgres{
		db: db,
	}

	utils.InitLogClient()

	transactionID := "testTransactionID"
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{transactionID},
			}},
	}

	inputContact := models.ContactRequest{
		Email:       "test@example.com",
		PhoneNumber: "123456789",
	}
	mock.ExpectBegin()


	// Mock database response for FindAllContacts
	mockRows := sqlmock.NewRows([]string{"id", "phoneNumber", "email", "linkedId", "linkPrecedence", "createdAt", "updatedAt", "deletedAt"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM contacts WHERE email = $1 OR phoneNumber = $2`)).
		WithArgs(inputContact.Email, inputContact.PhoneNumber).
		WillReturnRows(mockRows)

	// Define the query and expected arguments
	query := `INSERT INTO contacts (phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	expectedArgs := []driver.Value{
		inputContact.PhoneNumber,
		inputContact.Email,
		nil,
		"primary",
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	}

	id := 0

	// Define the expected result from the database
	rows := sqlmock.NewRows([]string{"id"}).AddRow(&id)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expectedArgs...).WillReturnRows(rows)


	mock.ExpectCommit()
	
	result, dbErr := p.FindOrCreateContact(ctx, inputContact)
	
	assert.Nil(t, dbErr)
	assert.NotNil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}


func TestFindOrCreateContactError(t *testing.T) {
	// mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// postgres struct with the mock DB
	p := postgres{
		db: db,
	}

	utils.InitLogClient()

	transactionID := "testTransactionID"
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{transactionID},
			}},
	}
	mock.ExpectBegin()

	inputContact := models.ContactRequest{
		Email:       "test@example.com",
		PhoneNumber: "123456789",
	}

	
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM contacts WHERE email = $1 OR phoneNumber = $2`)).
		WithArgs(inputContact.Email, inputContact.PhoneNumber).
		WillReturnError(fmt.Errorf("database error"))

	mock.ExpectCommit()
	_, dbErr := p.FindOrCreateContact(ctx, inputContact)

	
	assert.NotNil(t, dbErr)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

