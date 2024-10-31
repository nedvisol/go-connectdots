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

type CongressApiPagination struct {
	Count int     `json:"count"`
	Next  *string `json:"next"`
}

type CongressApiRequest struct {
	ContentType string `json:"contentType"`
	Format      string `json:"format"`
}

type CongressApiMemberResponse struct {
	Members    []*CongressApiMember   `json:"members"`
	Pagination *CongressApiPagination `json:"pagination"`
	Request    CongressApiRequest     `json:"request"`
}

type CongressApiCongress struct {
	Congresses []*CongressApiCongressItem `json:"congresses"`
	Pagination *CongressApiPagination     `json:"pagination"`
	Request    *CongressApiRequest        `json:"request"`
}

type CongressApiCongressItem struct {
	EndYear   *string               `json:"endYear"`
	Name      *string               `json:"name"`
	Sessions  []*CongressApiSession `json:"sessions"`
	StartYear *string               `json:"startYear"`
	URL       *string               `json:"url"`
}

type CongressApiSession struct {
	Chamber   *string `json:"chamber"`
	EndDate   *string `json:"endDate,omitempty"`
	Number    int     `json:"number"`
	StartDate *string `json:"startDate"`
	Type      *string `json:"type"`
}

type CongressApiBillsData struct {
	Bills      []*CongressApiBill     `json:"bills"`
	Pagination *CongressApiPagination `json:"pagination"`
	Request    *CongressApiRequest    `json:"request"`
}

type CongressApiBill struct {
	Congress                int                `json:"congress"`
	LatestAction            *CongressApiAction `json:"latestAction"`
	Number                  *string            `json:"number"`
	OriginChamber           *string            `json:"originChamber"`
	OriginChamberCode       *string            `json:"originChamberCode"`
	Title                   *string            `json:"title"`
	Type                    *string            `json:"type"`
	UpdateDate              *string            `json:"updateDate"`
	UpdateDateIncludingText *string            `json:"updateDateIncludingText"`
	URL                     *string            `json:"url"`
}

type CongressApiAction struct {
	ActionDate *string `json:"actionDate"`
	Text       *string `json:"text"`
}
