package processor

import "net/http"

type AuthenticationClient struct {
	scheme   string
	hostname string
}

func NewAuthenticationClient(scheme string, hostname string) *AuthenticationClient {
	return &AuthenticationClient{
		scheme:   scheme,
		hostname: hostname,
	}
}

func (this *AuthenticationClient) Do(*http.Request) (*http.Response, error) {
	panic("implement me")
}
