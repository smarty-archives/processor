package processor

import (
	"testing"

	"net/http"

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
	}
	this.verifier.Verify(input)

	this.AssertEqual("GET", this.client.request.Method)
	this.AssertEqual("/street-address?street=Street1", this.client.request.URL.String())
}

func (this *VerifierFixture) rawQuery() string {
	return this.client.request.URL.RawQuery
}

///////////////////////////////////////////////////////////////

type FakeHTTPClient struct {
	request *http.Request
}

func (this *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	this.request = request
	return nil, nil
}
