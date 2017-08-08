package processor

import "net/http"

type AuthenticationClient struct {
	inner    HTTPClient
	scheme   string
	hostname string
}

func NewAuthenticationClient(inner HTTPClient, scheme string, hostname string, authId string, authToken string) *AuthenticationClient {
	return &AuthenticationClient{
		inner:    inner,
		scheme:   scheme,
		hostname: hostname,
	}
}

func (this *AuthenticationClient) Do(request *http.Request) (*http.Response, error) {
	request.URL.Scheme = this.scheme
	request.URL.Host = this.hostname
	request.Host = this.hostname
	return this.inner.Do(request)
}
