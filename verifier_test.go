package processor

import (
	"testing"

	"net/http"

	"bytes"
	"io/ioutil"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestVerifierFixture(t *testing.T) {
	gunit.Run(new(VerifierFixture), t)
}

type VerifierFixture struct {
	*gunit.Fixture

	client   *FakeHTTPClient
	verifier *SmartyVerifier
}

func (this *VerifierFixture) Setup() {
	this.client = &FakeHTTPClient{}
	this.verifier = NewSmartyVerifier(this.client)
}
func NewSmartyVerifier(client HTTPClient) *SmartyVerifier {
	return &SmartyVerifier{
		client: client,
	}
}

func (this *VerifierFixture) TestRequestComposedProperly() {
	input := AddressInput{
		Street1: "Street1",
		City:    "City",
		State:   "State",
		ZIPCode: "ZIPCode",
	}

	this.client.Configure("[{}]", http.StatusOK, nil)

	this.verifier.Verify(input)

	this.So(this.client.request.Method, should.Equal, "GET")
	this.So(this.client.request.URL.Path, should.Equal, "/street-address")
	this.AssertQueryStringValue("street", input.Street1)
	this.AssertQueryStringValue("city", input.City)
	this.AssertQueryStringValue("state", input.State)
	this.AssertQueryStringValue("zipcode", input.ZIPCode)
}
func (this *VerifierFixture) AssertQueryStringValue(key, expected string) {
	query := this.client.request.URL.Query()
	this.So(query.Get(key), should.Equal, expected)
}

func (this *VerifierFixture) TestResponseParsed() {
	this.client.Configure(rawJSONOutput, http.StatusOK, nil)
	result := this.verifier.Verify(AddressInput{})
	this.So(result.DeliveryLine1, should.Equal, "1 Santa Claus Ln")
	this.So(result.LastLine, should.Equal, "North Pole AK 99705-9901")
}

const rawJSONOutput = `
[
	{
        "delivery_line_1": "1 Santa Claus Ln",
        "last_line": "North Pole AK 99705-9901"
    }
]`

///////////////////////////////////////////////////////////////

type FakeHTTPClient struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (this *FakeHTTPClient) Configure(responseText string, statusCode int, err error) {
	this.response = &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseText)),
		StatusCode: statusCode,
	}
	this.err = err
}

func (this *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	this.request = request
	return this.response, this.err
}
