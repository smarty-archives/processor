package processor

import (
	"net/http"
)

type AuthenticationClient struct {
	inner     HTTPClient
	scheme    string
	hostname  string
	authId    string
	authToken string
}

func NewAuthenticationClient(inner HTTPClient, scheme string, hostname string, authId string, authToken string) *AuthenticationClient {
	return &AuthenticationClient{
		inner:     inner,
		scheme:    scheme,
		hostname:  hostname,
		authId:    authId,
		authToken: authToken,
	}
}

func (this *AuthenticationClient) Do(request *http.Request) (*http.Response, error) {
	request.URL.Scheme = this.scheme
	request.URL.Host = this.hostname
	query := request.URL.Query()
	query.Set("auth-id", this.authId)
	query.Set("auth-token", this.authToken)
	request.URL.RawQuery = query.Encode()
	request.Host = this.hostname
	return this.inner.Do(request)
}
