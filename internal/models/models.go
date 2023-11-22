package models

import (
	"time"

	"github.com/lib/pq"
)

type ContactRequest struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type ContactResponse struct {
	PrimaryContactID    int      `json:"primary_contact_id"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phone_numbers"`
	SecondaryContactIDs *[]int   `json:"secondary_contact_ids,omitempty"`
}

type Contact struct {
	ID             *int        `json:"id"`
	PhoneNumber    string      `json:"phone_number"`
	Email          string      `json:"email"`
	LinkedID       *int        `json:"linked_id"`
	LinkPrecedence string      `json:"link_precedence"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DeletedAt      pq.NullTime `json:"deleted_at,omitempty"`
}
