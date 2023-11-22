package db

import (
	"database/sql"
	"fmt"
	"log"
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
	//response := []models.ContactResponse{}
	existingContact, err := p.FindContact(ctx, inputContact)
	fmt.Println("Err : ", err)
	fmt.Println("existingContact : ", existingContact)

	if existingContact != nil {
		// Contact already exists, update or create secondary contact
		fmt.Println("Existing contact found")
		fmt.Println("existingContact : ", existingContact.Email)
		fmt.Println("existingContact : ", existingContact.PhoneNumber)
		updatedContact, err := p.handleExistingContact(ctx, *existingContact, inputContact)
		if err != nil {
			return models.ContactResponse{}, &identityreconciliationerror.IdentityReconciliationError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("unable to get the contact details, err %v", err),
				Trace:   txid,
			}
		}
		fmt.Println("updatedContact : ", updatedContact)
		// contactRespose := models.ContactResponse{}
		// contactRespose.Emails = []string{contact.Email}
		// contactRespose.PhoneNumbers = []string{contact.PhoneNumber}
		// contactRespose.PrimaryContactID = *contact.ID
		return models.ContactResponse{}, nil
		//return *updatedContact, nil
	} else {
		fmt.Println("Existing contact not found")
		// Create a new primary contact
		contact := models.Contact{
			Email:          inputContact.Email,
			PhoneNumber:    inputContact.PhoneNumber,
			LinkPrecedence: "primary",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		id, err := p.saveContact(ctx, contact)
		fmt.Println("Error : ", err)
		if err != nil {
			return models.ContactResponse{}, &identityreconciliationerror.IdentityReconciliationError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("unable to get the contact details, err %v", err),
				Trace:   txid,
			}
		}
		contact.ID = id
		contactRespose := models.ContactResponse{}
		contactRespose.Emails = []string{contact.Email}
		contactRespose.PhoneNumbers = []string{contact.PhoneNumber}
		contactRespose.PrimaryContactID = *contact.ID
		//contactRespose.SecondaryContactIDs = []int{&contact.LinkedID}
		//response = append(response, contactRespose)
		fmt.Println("Response : ", contactRespose)
		return contactRespose, nil
	}
}

func (p postgres) FindContact(ctx *gin.Context, inputContact models.ContactRequest) (*models.Contact, *identityreconciliationerror.IdentityReconciliationError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	query := "SELECT id, email, phoneNumber, linkedId, linkPrecedence, createdAt, updatedAt, deletedAt FROM contacts WHERE email=$1 OR phoneNumber=$2"
	row := p.db.QueryRow(query, inputContact.Email, inputContact.PhoneNumber)

	contact := models.Contact{}
	err := row.Scan(&contact.ID,
		&contact.Email,
		&contact.PhoneNumber,
		&contact.LinkedID,
		&contact.LinkPrecedence,
		&contact.CreatedAt,
		&contact.UpdatedAt,
		&contact.DeletedAt,
	)
	fmt.Println("Error : ", err)
	if err != nil && err != sql.ErrNoRows {
		// if err == sql.ErrNoRows {
		// 	// Handle case where no rows were found
		// 	return contact, &limitoffererror.CreditCardError{
		// 		Code:    http.StatusNotFound,
		// 		Message: "account not found",
		// 		Trace:   txid,
		// 	}
		// }

		utils.Logger.Error(fmt.Sprintf("error while scanning contact details from db, txid : %v, error: %v", txid, err))
		return &contact, &identityreconciliationerror.IdentityReconciliationError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("unable to get the contact details, err %v", err),
			Trace:   txid,
		}
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	utils.Logger.Info(fmt.Sprintf("successfully fetched contact details from db, txid : %v", txid))

	return &contact, nil
}

func (p postgres) handleExistingContact(ctx *gin.Context, existingContact models.Contact, inputContact models.ContactRequest) (*models.Contact, *identityreconciliationerror.IdentityReconciliationError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	// Check if the incoming request introduces new information
	fmt.Println("existingContact.Email : ", existingContact.Email)
	fmt.Println("inputContact.Email : ", inputContact.Email)
	fmt.Println()
	fmt.Println("existingContact.PhoneNumber : ", existingContact.PhoneNumber)
	fmt.Println("inputContact.PhoneNumber : ", inputContact.PhoneNumber)
	if (existingContact.Email != inputContact.Email) || (existingContact.PhoneNumber != inputContact.PhoneNumber) {
		// Create a new secondary contact entry
		fmt.Println("Creating a secondary contact")
		secondaryContact := models.Contact{
			Email:          inputContact.Email,
			PhoneNumber:    inputContact.PhoneNumber,
			LinkPrecedence: "secondary",
			CreatedAt:      existingContact.CreatedAt,
			UpdatedAt:      time.Now(),
			//LinkedID:       existingContact.ID,
		}

		// Save the new secondary contact
		linkedID, err := p.saveContact(ctx, secondaryContact)
		if err != nil {
			return nil, &identityreconciliationerror.IdentityReconciliationError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("unable to save the contact details, err %v", err),
				Trace:   txid,
			}
		}

		fmt.Println("secondary contact stored successfully")

		//
		// existingContact.LinkedID = secondaryContact.ID
		// if err := p.saveContact(ctx, existingContact); err != nil {
		// 	return nil, err
		// }

		// Update the existing contact to link to the new secondary contact
		fmt.Println("linkedID.ID* : ", linkedID)
		fmt.Println("linkedID.ID : ", *linkedID)
		err = p.updateContact(ctx, existingContact, *linkedID)
		if err != nil {
			return nil, err
		}

		return &secondaryContact, nil
	} else {
		fmt.Println("No New information to be added for the existing user")
	}

	// No new information introduced, return the existing contact as is
	return &existingContact, nil
}

func (p postgres) updateContact(ctx *gin.Context, contact models.Contact, linkedID int64) *identityreconciliationerror.IdentityReconciliationError {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	_, err := p.db.Exec("UPDATE contacts SET linkedId = $1 WHERE id = $2", linkedID, contact.ID)
	fmt.Println("err 3 ", err)
	if err != nil {
		log.Println("error updating limit offer status:", err)
		return &identityreconciliationerror.IdentityReconciliationError{
			Code:    http.StatusInternalServerError,
			Message: "error while updating the linkedId of existing contact",
			Trace:   txid,
		}
	}
	fmt.Println("Updation is successful")
	return nil
}

func (p postgres) saveContact(ctx *gin.Context, contact models.Contact) (*int64, *identityreconciliationerror.IdentityReconciliationError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	if contact.ID == nil {
		// Insert new contact
		fmt.Println("Inserting a new row in saveContact")

		query := "INSERT INTO contacts (phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
		// result, err := p.db.Exec(query, contact.PhoneNumber, contact.Email, contact.LinkedID, contact.LinkPrecedence, contact.CreatedAt, contact.UpdatedAt)
		err := p.db.QueryRow(query, contact.PhoneNumber, contact.Email, contact.LinkedID, contact.LinkPrecedence, contact.CreatedAt, contact.UpdatedAt).Scan(&contact.ID)
		fmt.Println("Error while inserting a row : ", err)
		if err != nil {
			return nil, &identityreconciliationerror.IdentityReconciliationError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("unable to save the contact details, err %v", err),
				Trace:   txid,
			}
		}

		// lastInsertID, err := result.LastInsertId()
		// fmt.Println("Err : ", err , "lastInsertID : ", lastInsertID)
		// if err != nil {
		// 	return &identityreconciliationerror.IdentityReconciliationError{
		// 		Code:    http.StatusInternalServerError,
		// 		Message: fmt.Sprintf("unable to save the contact details, err %v", err),
		// 		Trace:   txid,
		// 	}
		// }
		// contact.ID = &lastInsertID
		fmt.Println("Last contactID is : ", *contact.ID)
	}
	// } else {

	// 	fmt.Println("Updating the existimg contact in saveContact")
	// 	// Update existing contact
	// 	query := "UPDATE contacts SET phoneNumber=$1, email=$2, linkedId=$3, linkPrecedence=$4, updatedAt=$4 WHERE id=$5"
	// 	_, err := p.db.Exec(query, contact.PhoneNumber, contact.Email, contact.LinkedID, contact.LinkPrecedence, contact.UpdatedAt, contact.ID)
	// 	if err != nil {
	// 		return &identityreconciliationerror.IdentityReconciliationError{
	// 			Code:    http.StatusInternalServerError,
	// 			Message: fmt.Sprintf("unable to save the contact details, err %v", err),
	// 			Trace:   txid,
	// 		}
	// 	}
	// }

	return contact.ID, nil
}
