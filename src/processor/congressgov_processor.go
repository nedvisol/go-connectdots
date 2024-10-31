package processor

import (
	"context"
	"encoding/json"
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

func (c *CongressGovProcessor) applyApiToken(url string) string {
	return fmt.Sprintf("%s&api_key=%s", url, c.apiToken)
}

func (c *CongressGovProcessor) createMemberNodeInfo(member *model.CongressApiMember) *graphdb.NodeInfo {
	names := strings.Split(member.Name, ", ")
	first, last := names[0], names[1]

	return &graphdb.NodeInfo{
		Id:    util.GetSHA512(fmt.Sprintf("%s-congress", member.BioguideID)),
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
		Id:    util.GetSHA512(fmt.Sprintf("bill/%s/%s", *bill.OriginChamberCode, *bill.Number)),
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

func (c *CongressGovProcessor) processCurrentMembers(data []byte) {
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
			c.ctx,
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
func (c *CongressGovProcessor) processBillActions(data []byte) {
	fmt.Printf("processing bill actions %d bytes\n", len(data))
}

func (c *CongressGovProcessor) processBills(data []byte) {
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

		c.dmgr.Download(
			c.ctx,
			downloadmgr.NewHttpGetRequest(c.applyApiToken(billActionsUrl)),
			c.processBillActions,
			downloadmgr.NewDownloadCacheOption(TEN_YEARS),
		)
		cnt++
	}
	fmt.Printf("updated %d bills\n", cnt)
}

var congressUrlRegex = regexp.MustCompile(`congress/(\d+)`)

func (c *CongressGovProcessor) processCongress(data []byte) {
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
				c.ctx,
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
