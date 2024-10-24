package processor

import (
	"context"
	"fmt"

	"github.com/nedvisol/go-connectdots/config"
	"github.com/nedvisol/go-connectdots/downloadmgr"
)

type CongressGovProcessor struct {
	ctx      context.Context
	dmgr     *downloadmgr.DownloadManager
	apiToken string
}

const MEMBERS_URL = "https://api.congress.gov/v3/member?format=json&currentMember=true&limit=250"

func (c *CongressGovProcessor) getCurrentMembersUrl() string {
	return fmt.Sprintf("%s&api_key=%s", MEMBERS_URL, c.apiToken)
}

func (c *CongressGovProcessor) processCurrentMembers(data []byte) {
	fmt.Printf("processing current memebers %d bytes", len(data))
}

// Start implements Processor.
func (c *CongressGovProcessor) Start() {
	c.dmgr.Download(
		c.ctx,
		downloadmgr.NewHttpGetRequest(c.getCurrentMembersUrl()),
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
