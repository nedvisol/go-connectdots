package processor

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/nedvisol/go-connectdots/config"
	"github.com/nedvisol/go-connectdots/downloadmgr"
	"github.com/nedvisol/go-connectdots/graphdb"
	"github.com/nedvisol/go-connectdots/model"
	"github.com/nedvisol/go-connectdots/util"
)

type key int

const billContextKey key = 0
const billActionContextKey key = 1

type CongressGovProcessor struct {
	ctx        context.Context
	dmgr       *downloadmgr.DownloadManager
	apiToken   string
	graphdbsvc graphdb.GraphDbService
}

const TEN_YEARS = time.Hour * 24 * 3650
const MEMBERS_URL = "https://api.congress.gov/v3/member?format=json&currentMember=true&limit=250"
const CONGRESS_URL = "https://api.congress.gov/v3/congress?format=json"
const BILLS_URL = "https://api.congress.gov//v3/bill/%s?format=json&limit=250"

func getPersonIdByBioguideId(bioguideId string) string {
	return util.GetSHA512(fmt.Sprintf("%s-congress", bioguideId))
}

func getBillIdByBillNumber(congress int, originChamberCode string, billNumber string) string {
	return util.GetSHA512(fmt.Sprintf("%d/bill/%s/%s", congress, originChamberCode, billNumber))
}

func getVotedEdgeIdByBillAction(bill *model.CongressApiBill, billAction *model.CongressApiBillAction) string {
	return util.GetSHA512(fmt.Sprintf("%d/bill/%s/%s/action/%s", bill.Congress, *bill.OriginChamberCode, *bill.Number, *billAction.ActionDate))
}

func (c *CongressGovProcessor) applyApiToken(url string) string {
	return fmt.Sprintf("%s&api_key=%s", url, c.apiToken)
}

func (c *CongressGovProcessor) createMemberNodeInfo(member *model.CongressApiMember) *graphdb.NodeInfo {
	names := strings.Split(member.Name, ", ")
	first, last := names[0], names[1]

	return &graphdb.NodeInfo{
		Id:    getPersonIdByBioguideId(member.BioguideID),
		Label: "Person",
		Attrs: &map[string]interface{}{
			"first":     first,
			"last":      last,
			"subtype":   "CongressMember",
			"party":     member.PartyName,
			"state":     member.State,
			"chamber":   member.Terms.Item[0].Chamber,
			"sourceUrl": member.URL,
		},
	}
}

func (c *CongressGovProcessor) createMember(member *model.CongressApiMember) {
	var err error
	personNode := c.createMemberNodeInfo(member)

	err = c.graphdbsvc.UpdateNode(personNode, true)

	if err != nil {
		panic(err)
	}
	fmt.Printf("member added/updated %s\n", member.Name)
}

func (c *CongressGovProcessor) createBillNode(bill *model.CongressApiBill) {
	var err error
	billNode := &graphdb.NodeInfo{
		Id:    getBillIdByBillNumber(bill.Congress, *bill.OriginChamberCode, *bill.Number),
		Label: "Bill",
		Attrs: &map[string]interface{}{
			"title":         bill.Title,
			"billType":      bill.Type,
			"originChamber": bill.OriginChamber,
			"congress":      bill.Congress,
			"url":           bill.URL,
		},
	}

	err = c.graphdbsvc.UpdateNode(billNode, true)

	if err != nil {
		panic(err)
	}
	fmt.Printf("bill added/updated %s\n", *bill.Number)
}

func (c *CongressGovProcessor) processCurrentMembers(ctx context.Context, data []byte) {
	fmt.Printf("processing current memebers %d bytes\n", len(data))

	var result model.CongressApiMemberResponse

	// Parse (unmarshal) the JSON into the map
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	if result.Pagination != nil && result.Pagination.Next != nil {
		fmt.Printf("found more things to download! %s\n", *result.Pagination.Next)
		c.dmgr.Download(
			ctx,
			downloadmgr.NewHttpGetRequest(c.applyApiToken(*result.Pagination.Next)),
			c.processCurrentMembers,
			downloadmgr.NewDownloadCacheOption(TEN_YEARS),
		)
	}

	var cnt = 0
	for _, member := range result.Members {
		c.createMember(member)
		cnt++
	}
	fmt.Printf("updated %d members\n", cnt)
}

func (c *CongressGovProcessor) createVotedForEdge(
	bill *model.CongressApiBill,
	billAction *model.CongressApiBillAction,
	bioguideId string,
	vote string,
) {
	votedFor := &graphdb.EdgeInfo{
		Label: "VOTED",
		Id:    getVotedEdgeIdByBillAction(bill, billAction),
		Left: &graphdb.NodeInfo{
			Id:    getPersonIdByBioguideId(bioguideId),
			Label: "Person",
		},
		Right: &graphdb.NodeInfo{
			Id:    getBillIdByBillNumber(bill.Congress, *bill.OriginChamberCode, *bill.Number),
			Label: "Bill",
		},
		Attrs: &map[string]interface{}{
			"vote": vote,
			"date": billAction.ActionDate,
		},
	}
	err := c.graphdbsvc.UpdateEdge(votedFor, true)
	if err != nil {
		panic(err)
	}
}

func (c *CongressGovProcessor) processHouseRollCallVote(ctx context.Context, data []byte) {
	fmt.Printf("processing house rollcall vote %d bytes\n", len(data))

	var result model.CongressApiHouseRollcallVote

	// Parse (unmarshal) the JSON into the map
	err := xml.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	//get bill object from context
	bill := ctx.Value(billContextKey)
	billAction := ctx.Value(billActionContextKey)

	for _, recordedVote := range result.VoteData.RecordedVotes {
		bioguideId := recordedVote.Legislator.NameID
		vote := recordedVote.Vote

		c.createVotedForEdge(bill.(*model.CongressApiBill), billAction.(*model.CongressApiBillAction), *bioguideId, *vote)
	}

}

func (c *CongressGovProcessor) processSenateRollCallVote(ctx context.Context, data []byte) {
	fmt.Printf("processing senate rollcall vote %d bytes\n", len(data))
}

func (c *CongressGovProcessor) processBillActions(ctx context.Context, data []byte) {
	fmt.Printf("processing bill actions %d bytes\n", len(data))

	var result model.CongressApiBillActionsPayload

	// Parse (unmarshal) the JSON into the map
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	var cnt = 0
	for _, action := range result.Actions {
		if action.RecordedVotes != nil && *action.Type == "Floor" {
			billActionCtx := context.WithValue(ctx, billActionContextKey, action)
			for _, recordedVote := range action.RecordedVotes {
				if recordedVote.URL != nil {
					if strings.Contains(*recordedVote.URL, "//clerk.house.gov") {
						c.dmgr.Download(
							billActionCtx,
							downloadmgr.NewHttpGetRequest(*recordedVote.URL),
							c.processHouseRollCallVote,
							downloadmgr.NewDownloadCacheOption(TEN_YEARS),
						)
					} else if strings.Contains(*recordedVote.URL, "//www.senate.gov") {
						c.dmgr.Download(
							billActionCtx,
							downloadmgr.NewHttpGetRequest(*recordedVote.URL),
							c.processSenateRollCallVote,
							downloadmgr.NewDownloadCacheOption(TEN_YEARS),
						)
					}
				}
			}
		}
		cnt++
	}
	fmt.Printf("updated %d bills\n", cnt)
}

func (c *CongressGovProcessor) processBills(ctx context.Context, data []byte) {
	fmt.Printf("processing bills %d bytes\n", len(data))

	var result model.CongressApiBillsData

	// Parse (unmarshal) the JSON into the map
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	var cnt = 0
	for _, bill := range result.Bills {
		c.createBillNode(bill)

		//download and process bills
		billActionsUrl := strings.ReplaceAll(*bill.URL, "?format=json", "/actions?format=json")

		//create new context with value
		billCtx := context.WithValue(ctx, billContextKey, bill)

		c.dmgr.Download(
			billCtx,
			downloadmgr.NewHttpGetRequest(c.applyApiToken(billActionsUrl)),
			c.processBillActions,
			downloadmgr.NewDownloadCacheOption(TEN_YEARS),
		)
		cnt++
	}
	fmt.Printf("updated %d bills\n", cnt)
}

var congressUrlRegex = regexp.MustCompile(`congress/(\d+)`)

func (c *CongressGovProcessor) processCongress(ctx context.Context, data []byte) {
	fmt.Printf("processing congress %d bytes\n", len(data))

	var result model.CongressApiCongress

	// Parse (unmarshal) the JSON into the map
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	var cnt = 0
	for _, congress := range result.Congresses {
		//c.createMember(member)
		if match := congressUrlRegex.FindStringSubmatch(*congress.URL); match != nil {
			congressNum := match[1]

			billsUrl := fmt.Sprintf(BILLS_URL, congressNum)

			//download and process bills
			c.dmgr.Download(
				ctx,
				downloadmgr.NewHttpGetRequest(c.applyApiToken(billsUrl)),
				c.processBills,
			)
		}
		cnt++
		if cnt > 2 {
			break
		}
	}
	fmt.Printf("updated %d congresses\n", cnt)
}

// Start implements Processor.
func (c *CongressGovProcessor) Start() {
	c.dmgr.Download(
		c.ctx,
		downloadmgr.NewHttpGetRequest(c.applyApiToken(MEMBERS_URL)),
		c.processCurrentMembers,
	)

	c.dmgr.Download(
		c.ctx,
		downloadmgr.NewHttpGetRequest(c.applyApiToken(CONGRESS_URL)),
		c.processCongress,
	)
}

func NewCongressGovProcessor(
	ctx context.Context,
	dmgr *downloadmgr.DownloadManager,
	config *config.Config,
	graphdbsvc graphdb.GraphDbService,
) *CongressGovProcessor {
	return &CongressGovProcessor{
		ctx:        ctx,
		dmgr:       dmgr,
		apiToken:   config.CongressGovToken,
		graphdbsvc: graphdbsvc,
	}
}
