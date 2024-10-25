package model

import (
	"time"
)

// Structs to represent the JSON data

type CongressApiMemberDepiction struct {
	Attribution string `json:"attribution"`
	ImageUrl    string `json:"imageUrl"`
}

type CongressApiMemberTermItem struct {
	Chamber   string `json:"chamber"`
	StartYear int    `json:"startYear"`
}

type CongressApiMemberTerms struct {
	Item []CongressApiMemberTermItem `json:"item"`
}

type CongressApiMember struct {
	BioguideID string                     `json:"bioguideId"`
	Depiction  CongressApiMemberDepiction `json:"depiction"`
	District   int                        `json:"district,omitempty"` // omitempty since some members don't have a district (e.g., Senators)
	Name       string                     `json:"name"`
	PartyName  string                     `json:"partyName"`
	State      string                     `json:"state"`
	Terms      CongressApiMemberTerms     `json:"terms"`
	UpdateDate time.Time                  `json:"updateDate"`
	URL        string                     `json:"url"`
}

type Pagination struct {
	Count int     `json:"count"`
	Next  *string `json:"next"`
}

type CongressApiMemberRequest struct {
	ContentType string `json:"contentType"`
	Format      string `json:"format"`
}

type CongressApiMemberResponse struct {
	Members    []*CongressApiMember     `json:"members"`
	Pagination *Pagination              `json:"pagination"`
	Request    CongressApiMemberRequest `json:"request"`
}
