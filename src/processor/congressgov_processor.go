package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nedvisol/go-connectdots/config"
	"github.com/nedvisol/go-connectdots/downloadmgr"
	"github.com/nedvisol/go-connectdots/model"
)

type CongressGovProcessor struct {
	ctx      context.Context
	dmgr     *downloadmgr.DownloadManager
	apiToken string
}

const MEMBERS_URL = "https://api.congress.gov/v3/member?format=json&currentMember=true&limit=250"

func (c *CongressGovProcessor) applyApiToken(url string) string {
	return fmt.Sprintf("%s&api_key=%s", url, c.apiToken)
}

func (c *CongressGovProcessor) processCurrentMembers(data []byte) {
	fmt.Printf("processing current memebers %d bytes", len(data))

	var result model.CongressMemberResponse

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
}

// Start implements Processor.
func (c *CongressGovProcessor) Start() {
	c.dmgr.Download(
		c.ctx,
		downloadmgr.NewHttpGetRequest(c.applyApiToken(MEMBERS_URL)),
		c.processCurrentMembers,
	)
}

func NewCongressGovProcessor(ctx context.Context, dmgr *downloadmgr.DownloadManager, config *config.Config) *CongressGovProcessor {
	return &CongressGovProcessor{
		ctx:      ctx,
		dmgr:     dmgr,
		apiToken: config.CongressGovToken,
	}
}
