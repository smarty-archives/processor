package processor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestPipelineFixture(t *testing.T) {
	gunit.Run(new(PipelineFixture), t)
}

type PipelineFixture struct {
	*gunit.Fixture

	reader   *ReadWriteSpyBuffer
	writer   *ReadWriteSpyBuffer
	client   *IntegrationHTTPClient
	pipeline *Pipeline
}

func (this *PipelineFixture) Setup() {
	this.reader = NewReadWriteSpyBuffer("")
	this.writer = NewReadWriteSpyBuffer("")
	this.client = &IntegrationHTTPClient{}
	this.pipeline = NewPipeline(ioutil.NopCloser(this.reader), this.writer, this.client, 2)
}

func (this *PipelineFixture) LongTestPipeline() {
	fmt.Fprintln(this.reader, "Street1,City,State,ZIPCode")
	fmt.Fprintln(this.reader, "A,B,C,D")
	fmt.Fprintln(this.reader, "A,B,C,D")

	err := this.pipeline.Process()

	this.So(this.writer.String(), should.Equal,
		"Status,DeliveryLine1,LastLine,City,State,ZIPCode\n"+
			"Deliverable,AA,BB,CC,DD,EE\n"+
			"Deliverable,AA,BB,CC,DD,EE\n")
	this.So(err, should.BeNil)
}

type IntegrationHTTPClient struct{}

func (this *IntegrationHTTPClient) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       NewReadWriteSpyBuffer(integrationJSONOutput),
		StatusCode: http.StatusOK,
	}, nil
}

const integrationJSONOutput = `
[
	{
        "delivery_line_1": "AA",
        "last_line": "BB",
        "components": {
            "city_name": "CC",
            "state_abbreviation": "DD",
            "zipcode": "EE"
        },
		"analysis": {
			"dpv_match_code": "Y",
			"dpv_vacant": "N",
			"active": "Y"
		}
    }
]`
