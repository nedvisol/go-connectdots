package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

const MEMBERS_URL = "https://api.congress.gov/v3/member?format=json&currentMember=true&limit=250"

func (c *CongressGovProcessor) applyApiToken(url string) string {
	return fmt.Sprintf("%s&api_key=%s", url, c.apiToken)
}

func (c *CongressGovProcessor) createNodeInfo(member *model.CongressApiMember) *graphdb.NodeInfo {
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
	personNode := c.createNodeInfo(member)
	if err := c.graphdbsvc.CreateNode(personNode); err != nil {
		panic(err)
	}

	err := c.graphdbsvc.CreateNode(personNode)
	if err != nil {
		panic(err)
	}
	fmt.Printf("member created %s\n", member.Name)
}

func (c *CongressGovProcessor) processCurrentMembers(data []byte) {
	fmt.Printf("processing current memebers %d bytes", len(data))

	var result model.CongressApiMemberResponse

	// Parse (unmarshal) the JSON into the map
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	if result.Pagination != nil && result.Pagination.Next != nil {
		fmt.Printf("found more things to download! %s", *result.Pagination.Next)
		c.dmgr.Download(
			c.ctx,
			downloadmgr.NewHttpGetRequest(c.applyApiToken(*result.Pagination.Next)),
			c.processCurrentMembers,
		)
	}

	for _, member := range result.Members {
		c.createMember(member)
	}
}

// Start implements Processor.
func (c *CongressGovProcessor) Start() {
	c.dmgr.Download(
		c.ctx,
		downloadmgr.NewHttpGetRequest(c.applyApiToken(MEMBERS_URL)),
		c.processCurrentMembers,
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
