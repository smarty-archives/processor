package processor

import "net/http"

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type SmartyVerifier struct {
	client HTTPClient
}

func (this *SmartyVerifier) Verify(AddressInput) AddressOutput {
	request, _ := http.NewRequest("GET", "", nil)
	this.client.Do(request)

	return AddressOutput{}
	//panic("implement me")
}
