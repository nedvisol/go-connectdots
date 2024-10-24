package model

import (
	"time"
)

// Structs to represent the JSON data

type Depiction struct {
	Attribution string `json:"attribution"`
	ImageUrl    string `json:"imageUrl"`
}

type TermItem struct {
	Chamber   string `json:"chamber"`
	StartYear int    `json:"startYear"`
}

type Terms struct {
	Item []TermItem `json:"item"`
}

type Member struct {
	BioguideID string    `json:"bioguideId"`
	Depiction  Depiction `json:"depiction"`
	District   int       `json:"district,omitempty"` // omitempty since some members don't have a district (e.g., Senators)
	Name       string    `json:"name"`
	PartyName  string    `json:"partyName"`
	State      string    `json:"state"`
	Terms      Terms     `json:"terms"`
	UpdateDate time.Time `json:"updateDate"`
	URL        string    `json:"url"`
}

type Pagination struct {
	Count int     `json:"count"`
	Next  *string `json:"next"`
}

type Request struct {
	ContentType string `json:"contentType"`
	Format      string `json:"format"`
}

type CongressMemberResponse struct {
	Members    []Member    `json:"members"`
	Pagination *Pagination `json:"pagination"`
	Request    Request     `json:"request"`
}
