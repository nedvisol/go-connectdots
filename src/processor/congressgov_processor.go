package processor

type CongressGovProcessor struct {
}

func (c *CongressGovProcessor) getMembers() {

}

// Start implements Processor.
func (c *CongressGovProcessor) Start(input interface{}) {
	panic("unimplemented")
}

func NewCongressGovProcessor() Processor {
	return &CongressGovProcessor{}
}
