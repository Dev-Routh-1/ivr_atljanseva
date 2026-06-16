package models

import "github.com/google/uuid"

type CitizenIVR struct {
	Language     string    `json:"language"`
	Pincode      string    `json:"pincode"`
	Ward         string    `json:"ward"`
	NagarsevakID uuid.UUID `json:"nagarsevak_id"`
}

type CitizenLookupResponse struct {
	Found bool `json:"found"`
	CitizenIVR
}

type RegisterCitizenRequest struct {
	PhoneNumber  string    `json:"phone_number" binding:"required"`
	Language     string    `json:"language" binding:"required"`
	Pincode      string    `json:"pincode" binding:"required"`
	Ward         string    `json:"ward" binding:"required"`
	NagarsevakID uuid.UUID `json:"nagarsevak_id" binding:"required"`
}