package processor

type VerifyHandler struct {
	input    chan *Envelope
	output   chan *Envelope
	verifier Verifier
}

type Verifier interface {
	Verify(AddressInput) AddressOutput
}

func NewVerifyHandler(in, out chan *Envelope, verifier Verifier) *VerifyHandler {
	return &VerifyHandler{
		input:    in,
		output:   out,
		verifier: verifier,
	}
}

func (this *VerifyHandler) Handle() {
	for envelope := range this.input {
		envelope.Output = this.verifier.Verify(envelope.Input)
		this.output <- envelope
	}
}
