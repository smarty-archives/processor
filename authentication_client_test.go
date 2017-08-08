package processor

import (
	"net/http/httptest"
	"testing"

	"net/http"

	"errors"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestAuthenticationClient(t *testing.T) {
	gunit.Run(new(AuthenticationClientFixture), t)
}

type AuthenticationClientFixture struct {
	*gunit.Fixture

	inner  *FakeHTTPClient
	client *AuthenticationClient
}

func (this *AuthenticationClientFixture) Setup() {
	this.inner = &FakeHTTPClient{}
	this.client = NewAuthenticationClient(this.inner, "http", "different-company.com", "authid", "authtoken")
}

func (this *AuthenticationClientFixture) TestProvideInformationAddedBeforeRequestIsSent() {
	request := httptest.NewRequest("GET", "/path", nil)

	this.client.Do(request)

	this.So(this.inner.request.URL.Scheme, should.Equal, "http")
	this.So(this.inner.request.Host, should.Equal, "different-company.com")
	this.So(this.inner.request.URL.Host, should.Equal, "different-company.com")
	this.So(this.inner.request.URL.Query().Get("auth-id"), should.Equal, "authid")
	this.So(this.inner.request.URL.Query().Get("auth-token"), should.Equal, "authtoken")
}

func (this *AuthenticationClientFixture) TestResponseAndErrorFromInnerClientReturned() {
	this.inner.response = &http.Response{StatusCode: http.StatusTeapot}
	this.inner.err = errors.New("HTTP Error")
	request := httptest.NewRequest("GET", "/path", nil)
	response, err := this.client.Do(request)

	this.So(response.StatusCode, should.Equal, http.StatusTeapot)
	this.So(err.Error(), should.Equal, "HTTP Error")
}
