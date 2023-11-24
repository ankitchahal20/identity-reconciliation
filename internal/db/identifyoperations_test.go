package db

import (
	"testing"
	"time"

	"github.com/identity-reconciliation/internal/models"
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

// func TestFindAllContacts(t *testing.T) {
// 	// Create a new mock database
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// Create an instance of the postgres struct with the mock DB
// 	p := postgres{
// 		db: db,
// 	}

// 	utils.InitLogClient()

// 	// Set up the test data
// 	transactionID := "testTransactionID"
// 	ctx := &gin.Context{
// 		Request: &http.Request{
// 			Header: http.Header{
// 				constants.TransactionID: []string{transactionID},
// 			}},
// 	}

// 	inputContact := models.ContactRequest{
// 		Email:       "test@example.com",
// 		PhoneNumber: "123456789",
// 	}

// 	// Mock database response for Query
// 	mockRows := sqlmock.NewRows([]string{"id", "phoneNumber", "email", "linkedId", "linkPrecedence", "createdAt", "updatedAt", "deletedAt"}).
// 		AddRow(1, "123456789", "test@example.com", nil, "primary", time.Now(), time.Now(), nil)

// 	query := "SELECT * FROM contacts WHERE email = $1 OR phoneNumber = $2"
// 	mock.ExpectQuery(query).
// 		WithArgs(inputContact.Email, inputContact.PhoneNumber).
// 		WillReturnRows(mockRows)

// 	fmt.Println("Expected Query:", query)

// 	// Call the function to test
// 	result, dbErr := p.findAllContacts(ctx, inputContact)
// 	fmt.Println("DBERR : ", dbErr)
// 	// Assert that the function behaves as expected
// 	assert.Nil(t, dbErr)
// 	assert.NotNil(t, result)
// 	assert.Len(t, result, 1)

// 	// Assert that the expected database queries were called
// 	assert.NoError(t, mock.ExpectationsWereMet())

// 	// Set up the expected SQL query and result
// 	// mockRows := sqlmock.NewRows([]string{"id", "note"}).
// 	// 	AddRow(1, "Note 1")

// 	// mock.ExpectQuery(`SELECT id, note FROM notes`).
// 	// 	WillReturnRows(mockRows)

// 	// // Call the function being tested
// 	// notes, notesErr := p.GetNotes(ctx)

// 	// // Assert that the returned error is nil
// 	// assert.Nil(t, notesErr)

// 	// // Assert the expected number of notes
// 	// expectedNotes := []models.Notes{
// 	// 	{NoteId: "1", Note: "Note 1"},
// 	// }
// 	// assert.Equal(t, expectedNotes, notes)

// 	// // Assert that all expectations were met
// 	// err = mock.ExpectationsWereMet()
// 	// assert.Nil(t, err)

// }
