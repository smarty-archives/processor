package processor

type (
	Envelope struct {
		Input  AddressInput
		Output AddressOutput
	}

	AddressInput struct {
		Street1 string
		City    string
		State   string
		ZIPCode string
	}

	AddressOutput struct {
		DeliveryLine1 string
	}
)
