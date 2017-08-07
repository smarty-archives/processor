package processor

type VerifyHandler struct {
	input    chan *Envelope
	output   chan *Envelope
	verifier Verifier
}

type Verifier interface {
	Verify(AddressInput) AddressOutput
}

func NewVerifyHandler(input, output chan *Envelope, verifier Verifier) *VerifyHandler {
	return &VerifyHandler{
		input:    input,
		output:   output,
		verifier: verifier,
	}
}

func (this *VerifyHandler) Handle() {
	for envelope := range this.input {
		envelope.Output = this.verifier.Verify(envelope.Input)
		this.output <- envelope
	}
}
