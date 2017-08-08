package processor

import (
	"fmt"
	"testing"

	"net/http"

	"bytes"

	"errors"

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
	this.So(result.City, should.Equal, "North Pole")
	this.So(result.State, should.Equal, "AK")
	this.So(result.ZIPCode, should.Equal, "99705")
}

const rawJSONOutput = `
[
	{
        "delivery_line_1": "1 Santa Claus Ln",
        "last_line": "North Pole AK 99705-9901",
        "components": {
            "city_name": "North Pole",
            "state_abbreviation": "AK",
            "zipcode": "99705"
        }
    }
]`

func (this *VerifierFixture) TestMalformedJSONHandled() {
	const malformedRawJSONOutput = `alert('Hello, World!' DROP TABLE Users);`
	this.client.Configure(malformedRawJSONOutput, http.StatusOK, nil)
	result := this.verifier.Verify(AddressInput{})
	this.So(result.Status, should.Equal, "Invalid API Response")
}

func (this *VerifierFixture) TestHTTPErrorHandled() {
	this.client.Configure("", 0, errors.New("GOPHERS!"))

	result := this.verifier.Verify(AddressInput{})

	this.So(result.Status, should.Equal, "Invalid API Response")
}

func (this *VerifierFixture) TestHTTPResponseBodyClosed() {
	this.client.Configure(rawJSONOutput, http.StatusOK, nil)
	this.verifier.Verify(AddressInput{})
	this.So(this.client.responseBody.closed, should.Equal, 1)
}

func (this *VerifierFixture) TestDeliverableAddressStatus() {
	this.client.Configure(buildAnalysisJSON("Y", "N", "Y"), http.StatusOK, nil)
	output := this.verifier.Verify(AddressInput{})
	this.So(output.Status, should.Equal, "Deliverable")
}

func (this *VerifierFixture) TestValidVacantAddress() {
	this.client.Configure(buildAnalysisJSON("Y", "Y", "Y"), http.StatusOK, nil)
	output := this.verifier.Verify(AddressInput{})
	this.So(output.Status, should.Equal, "Vacant")
}

func (this *VerifierFixture) TestValidInactiveAddress() {
	this.client.Configure(buildAnalysisJSON("Y", "N", "?"), http.StatusOK, nil)
	output := this.verifier.Verify(AddressInput{})
	this.So(output.Status, should.Equal, "Inactive")
}

func buildAnalysisJSON(match, vacant, active string) string {
	template := `
	[
		{
			"analysis": {
				"dpv_match_code": "%s",
				"dpv_vacant": "%s",
				"active": "%s"
			}
		}
	]`
	return fmt.Sprintf(template, match, vacant, active)
}

///////////////////////////////////////////////////////////////

type FakeHTTPClient struct {
	request      *http.Request
	response     *http.Response
	responseBody *SpyBuffer
	err          error
}

func (this *FakeHTTPClient) Configure(responseText string, statusCode int, err error) {
	if err == nil {
		this.responseBody = NewSpyBuffer(responseText)
		this.response = &http.Response{
			Body:       this.responseBody,
			StatusCode: statusCode,
		}
	}
	this.err = err
}
func (this *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	this.request = request
	return this.response, this.err
}

///////////////////////////////////////////////////////////////

type SpyBuffer struct {
	*bytes.Buffer
	closed int
}

func NewSpyBuffer(value string) *SpyBuffer {
	return &SpyBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (this *SpyBuffer) Close() error {
	this.closed++
	this.Buffer.Reset()
	return nil
}
