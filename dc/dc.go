package dc

import "encoding/xml"

type dcResult struct {
	XMLName xml.Name    `xml:"OAI-PMH"`
	Records []DCLinks   `xml:"ListRecords>record"`
	Error   OAIPMHError `xml:"error"`
	RT      string      `xml:"ListRecords>resumptionToken"`
}

// DCLinks represents the XML data for a set of links under an individual record
type DCLinks struct {
	ContentURLS []string `xml:"metadata>dc>identifier"`
}

// OAIPMHError holds any errors details returned in the body. OAI_PMH uses this to report
// faults with the query rather than relying on HTTP Status
type OAIPMHError struct {
	Code        string `xml:"code,attr"`
	Description string `xml:",chardata"`
}

type DC struct {
	fetcher fetcher.Fetcher
}

func NewDC(f fetcher.Fetcher) *DC {
	return &DC{fetcher: f}
}

// Process will read data from the specified source, following nay resumption tokens and will add
// any ContentURLS found and add them to the channel.
func (dc *DC) Process(source string, links chan<- string) error
	return nil
}