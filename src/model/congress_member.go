package model

import (
	"encoding/xml"
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

type CongressApiBillActionsPayload struct {
	Actions    []*CongressApiBillAction `json:"actions"`
	Pagination *CongressApiPagination   `json:"pagination"`
	Request    *CongressApiRequest      `json:"request"`
}

type CongressApiBillAction struct {
	ActionCode    *string                    `json:"actionCode,omitempty"`
	ActionDate    *string                    `json:"actionDate"`
	SourceSystem  *CongressApiSourceSystem   `json:"sourceSystem"`
	Text          *string                    `json:"text"`
	Type          *string                    `json:"type"`
	RecordedVotes []*CongressApiRecordedVote `json:"recordedVotes,omitempty"`
}

type CongressApiSourceSystem struct {
	Code *int    `json:"code,omitempty"`
	Name *string `json:"name"`
}

type CongressApiRecordedVote struct {
	Chamber       *string   `json:"chamber"`
	Congress      *int      `json:"congress"`
	Date          time.Time `json:"date"`
	RollNumber    *int      `json:"rollNumber"`
	SessionNumber *int      `json:"sessionNumber"`
	URL           *string   `json:"url"`
}

type CongressApiHouseRollcallVote struct {
	XMLName      xml.Name                      `xml:"rollcall-vote"`
	VoteMetadata *CongressApiHouseVoteMetadata `xml:"vote-metadata"`
	VoteData     *CongressApiHouseVoteData     `xml:"vote-data"`
}

type CongressApiHouseVoteMetadata struct {
	Majority     *string                     `xml:"majority"`
	Congress     *string                     `xml:"congress"`
	Session      *string                     `xml:"session"`
	Chamber      *string                     `xml:"chamber"`
	RollcallNum  *string                     `xml:"rollcall-num"`
	LegisNum     *string                     `xml:"legis-num"`
	VoteQuestion *string                     `xml:"vote-question"`
	VoteType     *string                     `xml:"vote-type"`
	VoteResult   *string                     `xml:"vote-result"`
	ActionDate   *string                     `xml:"action-date"`
	ActionTime   *CongressApiHouseActionTime `xml:"action-time"`
	VoteDesc     *string                     `xml:"vote-desc"`
	VoteTotals   *CongressApiHouseVoteTotals `xml:"vote-totals"`
}

type CongressApiHouseActionTime struct {
	TimeETZ *string `xml:"time-etz,attr"`
	Text    *string `xml:",chardata"`
}

type CongressApiHouseVoteTotals struct {
	TotalsByPartyHeader *CongressApiHouseTotalsByPartyHeader `xml:"totals-by-party-header"`
	TotalsByParty       []*CongressApiHouseTotalsByParty     `xml:"totals-by-party"`
	TotalsByVote        *CongressApiHouseTotalsByVote        `xml:"totals-by-vote"`
}

type CongressApiHouseTotalsByPartyHeader struct {
	PartyHeader     *string `xml:"party-header"`
	YeaHeader       *string `xml:"yea-header"`
	NayHeader       *string `xml:"nay-header"`
	PresentHeader   *string `xml:"present-header"`
	NotVotingHeader *string `xml:"not-voting-header"`
}

type CongressApiHouseTotalsByParty struct {
	Party          *string `xml:"party"`
	YeaTotal       *string `xml:"yea-total"`
	NayTotal       *string `xml:"nay-total"`
	PresentTotal   *string `xml:"present-total"`
	NotVotingTotal *string `xml:"not-voting-total"`
}

type CongressApiHouseTotalsByVote struct {
	TotalStub      *string `xml:"total-stub"`
	YeaTotal       *string `xml:"yea-total"`
	NayTotal       *string `xml:"nay-total"`
	PresentTotal   *string `xml:"present-total"`
	NotVotingTotal *string `xml:"not-voting-total"`
}

type CongressApiHouseVoteData struct {
	RecordedVotes []*CongressApiHouseRecordedVote `xml:"recorded-vote"`
}

type CongressApiHouseRecordedVote struct {
	Legislator *CongressApiHouseLegislator `xml:"legislator"`
	Vote       *string                     `xml:"vote"`
}

type CongressApiHouseLegislator struct {
	NameID         *string `xml:"name-id,attr"`
	SortField      *string `xml:"sort-field,attr"`
	UnaccentedName *string `xml:"unaccented-name,attr"`
	Party          *string `xml:"party,attr"`
	State          *string `xml:"state,attr"`
	Role           *string `xml:"role,attr"`
	Text           *string `xml:",chardata"`
}
