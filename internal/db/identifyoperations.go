package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	identityreconciliationerror "github.com/identity-reconciliation/internal/IdentityReconciliationError"
	"github.com/identity-reconciliation/internal/constants"
	"github.com/identity-reconciliation/internal/models"
	"github.com/identity-reconciliation/internal/utils"
)

func (p postgres) FindOrCreateContact(ctx *gin.Context, inputContact models.ContactRequest) (models.ContactResponse, *identityreconciliationerror.IdentityReconciliationError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	
	tx, dbErr := p.db.Begin()
    if dbErr != nil {
        utils.Logger.Error(fmt.Sprintf("error starting database transaction, txid: %v, err: %v", txid, dbErr))
        return models.ContactResponse{}, &identityreconciliationerror.IdentityReconciliationError{
            Code:    http.StatusInternalServerError,
            Message: fmt.Sprintf("unable to start a database transaction, err: %v", dbErr),
            Trace:   txid,
        }
    }

	contacts, err := p.findAllContacts(tx, ctx, inputContact)
	if err != nil {
		tx.Rollback()
		utils.Logger.Error(fmt.Sprintf("error while reteriving all contacts, txid : %v", txid))
		return models.ContactResponse{}, err
	}

	var contactList []models.Contact
	if len(contacts) != 0 {
		utils.Logger.Info(fmt.Sprintf("existing contact found for the given request, txid : %v", txid))
		contactResponse, err :=  p.handleExistingContact(tx, ctx, contacts, inputContact)
		if err != nil {
			tx.Rollback()
			return models.ContactResponse{}, err
		}
		return contactResponse, nil
	}

	utils.Logger.Info(fmt.Sprintf("no existing contact found for the given request, txid : %v", txid))

	// Create a new entry if no existing contacts match
	newContact := models.Contact{
		Email:          inputContact.Email,
		PhoneNumber:    inputContact.PhoneNumber,
		LinkPrecedence: "primary",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := "INSERT INTO contacts (phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	dbErr = tx.QueryRow(query, newContact.PhoneNumber, newContact.Email, newContact.LinkedID, newContact.LinkPrecedence, newContact.CreatedAt, newContact.UpdatedAt).Scan(&newContact.ID)
	if dbErr != nil {
		tx.Rollback()
		utils.Logger.Error(fmt.Sprintf("error while creating a primary contact, txid : %v", txid))
		return models.ContactResponse{}, &identityreconciliationerror.IdentityReconciliationError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("unable to fetch the contact details, err %v", dbErr),
			Trace:   txid,
		}
	}
	utils.Logger.Info(fmt.Sprintf("new contact created for the given request, txid : %v", txid))
	// Commit the transaction if everything is successful
    if err := tx.Commit(); err != nil {
        utils.Logger.Error(fmt.Sprintf("error committing database transaction, txid: %v, err: %v", txid, err))
        return models.ContactResponse{}, &identityreconciliationerror.IdentityReconciliationError{
            Code:    http.StatusInternalServerError,
            Message: fmt.Sprintf("unable to commit the database transaction, err: %v", err),
            Trace:   txid,
        }
    }
	contactList = append(contactList, newContact)
	return transformContact(contactList), nil
}

func (p postgres) handleExistingContact(tx *sql.Tx, ctx *gin.Context, contacts []models.Contact, inputContact models.ContactRequest) (models.ContactResponse, *identityreconciliationerror.IdentityReconciliationError) {

	if len(contacts) == 1 {
		if contacts[0].Email != inputContact.Email || contacts[0].PhoneNumber != inputContact.PhoneNumber {
			contactList, err := p.foundOneRecord(tx, ctx, contacts, inputContact.Email, inputContact.PhoneNumber)
			if err != nil {
				return models.ContactResponse{}, err
			}
			return transformContact(contactList), nil
		} else {
			return transformContact(contacts), nil
		}
	}

	if len(contacts) > 1 {
		contactList, err := p.foundMultipleRecord(tx, ctx, contacts, inputContact.Email, inputContact.PhoneNumber)
		if err != nil {
			return models.ContactResponse{}, err
		}
		return transformContact(contactList), nil
	}

	return models.ContactResponse{}, nil
}

func (p postgres) foundOneRecord(tx *sql.Tx, ctx *gin.Context, contacts []models.Contact, email string, phoneNumber string) ([]models.Contact, *identityreconciliationerror.IdentityReconciliationError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	oldRecord := contacts[0]

	// If by chance there is a problematic record, restore consistency
	if oldRecord.LinkPrecedence == "secondary" {
		oldRecord.LinkPrecedence = "primary"
		_, err := tx.Exec("UPDATE contacts SET linkPrecedence = $1 WHERE id = $2", oldRecord.LinkPrecedence, oldRecord.ID)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("error while updating the contacts info, txid : %v", txid))
			return nil, &identityreconciliationerror.IdentityReconciliationError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("unable to update the contact details, err %v", err),
				Trace:   txid,
			}
		}
	}

	newRecord := models.Contact{
		Email:          email,
		PhoneNumber:    phoneNumber,
		LinkedID:       oldRecord.ID,
		LinkPrecedence: "secondary",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	query := "INSERT INTO contacts (phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err := tx.QueryRow(query, newRecord.PhoneNumber, newRecord.Email, newRecord.LinkedID, newRecord.LinkPrecedence, newRecord.CreatedAt, newRecord.UpdatedAt).Scan(&newRecord.ID)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("error while creating a secondary contact, txid : %v", txid))
		return nil, &identityreconciliationerror.IdentityReconciliationError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("unable to add new contact details, err %v", err),
			Trace:   txid,
		}
	}

	return []models.Contact{oldRecord, newRecord}, nil
}

// foundMultipleRecord function
func (p postgres) foundMultipleRecord(tx *sql.Tx, ctx *gin.Context, contacts []models.Contact, email, phoneNumber string) ([]models.Contact, *identityreconciliationerror.IdentityReconciliationError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	primaryRec := contacts[0]

	var statements []string
	statements = append(statements, fmt.Sprintf("UPDATE contacts SET linkPrecedence = 'primary' WHERE id = %d", *primaryRec.ID))

	for index, contact := range contacts {
		if index > 0 && (contact.LinkPrecedence != "secondary" || *contact.LinkedID != *primaryRec.ID) {
			statements = append(statements, fmt.Sprintf("UPDATE contacts SET linkPrecedence = 'secondary', linkedId = %d WHERE id = %d", *primaryRec.ID, *contact.ID))
		}
	}

	// Execute SQL update statements
	for _, statement := range statements {
		_, err := tx.Exec(statement)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("error while updating contacts information, txid : %v", txid))
			return nil, &identityreconciliationerror.IdentityReconciliationError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("unable to add new contact details, err %v", err),
				Trace:   txid,
			}
		}
	}

	newRecord := models.Contact{
		Email:          email,
		PhoneNumber:    phoneNumber,
		LinkedID:       primaryRec.ID,
		LinkPrecedence: "secondary",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := "INSERT INTO contacts (phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err := tx.QueryRow(query, newRecord.PhoneNumber, newRecord.Email, newRecord.LinkedID, newRecord.LinkPrecedence, newRecord.CreatedAt, newRecord.UpdatedAt).Scan(&newRecord.ID)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("error while creating a secondary contact, txid : %v", txid))
		return nil, &identityreconciliationerror.IdentityReconciliationError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("unable to add new contact details, err %v", err),
			Trace:   txid,
		}
	}

	return append(contacts, newRecord), nil
}

func (p postgres) findAllContacts(tx *sql.Tx, ctx *gin.Context, inputContact models.ContactRequest) ([]models.Contact, *identityreconciliationerror.IdentityReconciliationError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	var contacts []models.Contact
	query := "SELECT * FROM contacts WHERE email = $1 OR phoneNumber = $2"
	rows, err := tx.Query(query, inputContact.Email, inputContact.PhoneNumber)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("error while fetching all contacts information, txid : %v", txid))
		return []models.Contact{}, &identityreconciliationerror.IdentityReconciliationError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("unable to fetch the contact details, err %v", err),
			Trace:   txid,
		}
	}
	defer rows.Close()

	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(
			&contact.ID,
			&contact.PhoneNumber,
			&contact.Email,
			&contact.LinkedID,
			&contact.LinkPrecedence,
			&contact.CreatedAt,
			&contact.UpdatedAt,
			&contact.DeletedAt,
		); err != nil {
			utils.Logger.Error(fmt.Sprintf("error while scanning the fetched records, txid : %v", txid))
			return nil, &identityreconciliationerror.IdentityReconciliationError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("unable to scan the contact details, err %v", err),
				Trace:   txid,
			}
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

// transformContact function
func transformContact(contacts []models.Contact) models.ContactResponse {
	var (
		emails              []string
		phoneNumbers        []string
		secondaryContactIds []int
		primaryContactID    int
	)
	// this function will only be called when atleast one contact is found in db for the given request,
	if len(contacts) != 0 {
		primaryContactID = *contacts[0].ID
	}

	for index, contact := range contacts {
		if contact.Email != "" && !utils.Contains(emails, contact.Email) {
			emails = append(emails, contact.Email)
		}
		if contact.PhoneNumber != "" && !utils.Contains(phoneNumbers, contact.PhoneNumber) {
			phoneNumbers = append(phoneNumbers, contact.PhoneNumber)
		}
		if index > 0 {
			secondaryContactIds = append(secondaryContactIds, int(*contact.ID))
		}
	}

	contactResponse := models.ContactResponse{}
	contactResponse.Emails = emails
	contactResponse.PhoneNumbers = phoneNumbers
	contactResponse.PrimaryContactID = primaryContactID
	if len(secondaryContactIds) != 0 {
		contactResponse.SecondaryContactIDs = &secondaryContactIds
	}
	return contactResponse
}
