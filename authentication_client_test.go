package processor

import (
	"testing"

	"github.com/smartystreets/gunit"
	"net/http/httptest"
	"github.com/smartystreets/assertions/should"
)

func TestAuthenticationClient(t *testing.T) {
	gunit.Run(new(AuthenticationClientFixture), t)
}


type AuthenticationClientFixture struct {
	*gunit.Fixture

	inner *FakeHTTPClient
	client *AuthenticationClient
}
func (this *AuthenticationClientFixture) Setup() {
	this.inner = &FakeHTTPClient{}
	this.client = &AuthenticationClient{}
}

func (this *AuthenticationClientFixture) TestHostnameAndScheme() {
	this.client.WithHostname("us-street.api.smartystreets.com")
	request := httptest.NewRequest("GET", "/path", nil)

	this.client.Do(request)

	this.So(this.inner.request.Host, should.Equal, "us-street.api.smartystreets.com")
	this.So(this.inner.request.URL.Scheme, should.Equal, "https")
}

